var async = require("async"),
    path = require("path"),
	aws = require("aws-sdk"),
    fs = require("fs"),
	ec2 = require("./ec2.js");

var config = JSON.parse(fs.readFileSync("config.json"));
var awsConfig = JSON.parse(fs.readFileSync("aws.json"));

aws.config.update({accessKeyId: awsConfig.awsKey, secretAccessKey: awsConfig.awsSecret});
aws.config.update({region: 'ap-southeast-2'});

// upload to S3
async.series({
	/*deleteOldExe : function(cb) {
			console.log("Deleting old version of the executable from S3...");
			
			// first we delete the old file
			new aws.S3().deleteObject({Bucket: awsConfig.s3Bucket, Key: awsConfig.s3Key}, function(err, data) {
				if(err)
					console.log("Unable to delete", awsConfig.s3Bucket+ "/"+ awsConfig.s3Key, "because:", err);
				else
					console.log("Previous version of", awsConfig.s3Bucket+ "/"+ awsConfig.s3Key, "removed");
				cb();
			});
		},
	uploadNewExe: function(cb) {
			console.log("Uploading the new executable...");
			
			if(!fs.existsSync(config.logbuzzExecutable))
				throw new Exception('File', srcFile, "does not exist, unable to continue")
			
			var options = {
				Bucket: awsConfig.s3Bucket,
				Key: awsConfig.s3Key,
				ACL: "public-read",
				Body: fs.readFileSync(config.logbuzzExecutable)
			}
			
			new aws.S3().putObject(options, function(err, data) {
			  
				if (err != null)
					throw new Exception("Error uploading file:", config.logbuzzExecutable, ":", err)
				
				cb();
			});
			
		},*/
	shutdownOldTestInstance: function(cb) {
			console.log("Shutting down all old log test instances");
			
			ec2.getRunningInstancesByGroup(awsConfig.ec2GroupName, function(err, data) {
				if(err != null)
					throw new Exception(err);
				
				console.log("We have", data.length, "test instances running");
				
				var terminatedCount = 0;
				
				for(i = 0; i < data.length; i++) {
					console.log("About to shut down instance: " + data[i].instanceID);
					ec2.terminateInstance(data[i].instanceID, function(err, data) {
						if(err != null)
							throw new Exception(err);
						
						terminateCount++;
						
						if(terminatedCount == data.length)
							cb();
					});
				}
				
				
			});
			
		},
	startNewLogInstance: function(cb) {
			
			//http://fratuz610-deploy.s3.amazonaws.com/logbuzz/logbuzz
			var userDataList = [];
			userDataList.push("#!/bin/bash");
			userDataList.push("sudo groupadd logbuzz");
			userDataList.push("sudo useradd -g logbuzz -d /opt/logbuzz -m -s /bin/bash logbuzz");
			userDataList.push("sudo su logbuzz");
			userDataList.push("cd ~");
			userDataList.push("wget http://" + awsConfig.s3Bucket + ".s3.amazonaws.com/" + awsConfig.s3Key);
			userDataList.push("chmod u+x ./logbuzz");
			userDataList.push("nohup logbuzz > /dev/null 2>&1 &");
		
			console.log("Starting the new log instance");
			
			ec2.runInstances(1, userDataList, function(err, data) {
				if(err != null)
					throw new Exception(err);
				
				var startedInstanceID = data[0];
				
				console.log("Successfully started instance:", startedInstanceID);
				
				var instanceRunning = false;
				
				async.whilst(
					function () { return !instanceRunning },
					function (callback) {
						ec2.isInstanceRunning(startedInstanceID, function(err, data) {
							if(err != null)
								throw new Exception(err);
							
							console.log("Instance running ?", instanceRunning);
							
							instanceRunning = data.running;
							
							setTimeout(callback, 1000);
						});
						
					},
					function (err) {
						console.log("Instance", startedInstanceID, "running now");
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

