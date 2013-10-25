var async = require("async"),
    path = require("path"),
	aws = require("aws-sdk"),
    fs = require("fs"),
	os = require("os"),
	ec2 = require("./ec2.js");
	Zip = require("./zip.js");

var config = JSON.parse(fs.readFileSync("config.json"));
var awsConfig = JSON.parse(fs.readFileSync("aws.json"));

aws.config.update({accessKeyId: awsConfig.awsKey, secretAccessKey: awsConfig.awsSecret});
aws.config.update({region: 'ap-southeast-2'});

// upload to S3
async.series({
	
	deleteOldBundle: function(cb) {
		
		new aws.S3().deleteObject({Bucket: awsConfig.s3Bucket, Key: awsConfig.s3Key}, function(err, data) {
				if(err)
					console.log("Unable to delete", awsConfig.s3Bucket+ "/"+ awsConfig.s3Key, "because:", err);
				else
					console.log("Previous version of", awsConfig.s3Bucket+ "/"+ awsConfig.s3Key, "removed");
				cb();
			});
		
	},
	uploadNewBundle: function(cb) {
		// we get the static files (web)
		
		var zip = new Zip();
		zip.addFile(config.logbuzzExecutable, "logbuzz");
		zip.addFolder(config.logbuzzStaticFilesFolder, "static");
		var tempBundle = zip.toTempFile();
		
		console.log("Temporary bundle saved at", tempBundle, "uploading it");
		
		var options = {
			Bucket: awsConfig.s3Bucket,
			Key: awsConfig.s3Key,
			ACL: "public-read",
			Body: fs.readFileSync(tempBundle)
		}
		
		new aws.S3().putObject(options, function(err, data) {
		  
			if (err != null)
				throw new Error("Error uploading file:", config.logbuzzExecutable, ":", err)
			
			fs.unlinkSync(tempBundle);
			
			cb();
		});
		
	},
	shutdownOldTestInstance: function(cb) {
			
			console.log("Shutting down all old log test instances");
			
			ec2.getRunningInstancesByGroup(awsConfig.ec2GroupName, function(err, runningInstanceList) {
				if(err != null)
					throw new Error(err);
				
				console.log("We have", runningInstanceList.length, "test instances running");
				
				var terminatedCount = 0;
				
				if(runningInstanceList.length == 0) {
					cb();
					return;
				}
				
				for(i = 0; i < runningInstanceList.length; i++) {
					console.log("About to shut down instance: " + runningInstanceList[i].instanceID);
					ec2.terminateInstance(runningInstanceList[i].instanceID, function(err, data) {
						if(err != null)
							throw new Error(err);
						
						terminatedCount++;
						
						if(terminatedCount == runningInstanceList.length)
							cb();
					});
				}
				
				
			});
			
		},
	startNewLogInstance: function(cb) {
			
			var userDataList = [];
			userDataList.push("#!/bin/bash");
			userDataList.push("sudo apt-get install unzip");
			userDataList.push("sudo groupadd logbuzz");
			userDataList.push("sudo useradd -g logbuzz -d /opt/logbuzz -m -s /bin/bash logbuzz");
			userDataList.push("sudo iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 5222 -j REDIRECT --to-ports 1212");
			userDataList.push("sudo -u logbuzz wget http://" + awsConfig.s3Bucket + ".s3.amazonaws.com/" + awsConfig.s3Key +" -P /tmp");
			userDataList.push("sudo -u logbuzz unzip /tmp/logbuzz.zip -d /opt/logbuzz");
			userDataList.push("sudo -u logbuzz chmod u+x /opt/logbuzz/logbuzz");
			userDataList.push("sudo -u logbuzz nohup /opt/logbuzz/logbuzz > /dev/null 2>&1 &");
		
			console.log("Starting the new log instance");
			
			ec2.runInstances(1, userDataList, function(err, data) {
				if(err != null)
					throw new Error(err);
				
				var startedInstanceID = data[0];
				
				console.log("Successfully started instance:", startedInstanceID);
				
				var instanceRunning = false;
				var privateIpAddress = null;
				var publicIpAddress = null;
				
				console.log("Waiting for the instance to come online...");
				async.whilst(
					function () { return !instanceRunning },
					function (callback) {
						ec2.isInstanceRunning(startedInstanceID, function(err, data) {
							if(err != null)
								throw new Error(err);
							
							instanceRunning = data.running;
							
							if(data.privateIpAddress) privateIpAddress = data.privateIpAddress;
							if(data.publicIpAddress) publicIpAddress = data.publicIpAddress;
							
							setTimeout(callback, 1000);
						});
						
					},
					function (err) {
						console.log("Instance", startedInstanceID, "running now at", privateIpAddress, "and", publicIpAddress);
						cb(err);
					}
				);
				
				
			});
			
		}
	},
	function(err, result) {
		if(err)
			console.log("Deployment failed", err);
		else
			console.log("Deployment completed");
	}
);

