var	aws = require("aws-sdk"),
    fs = require("fs");

var awsConfig = JSON.parse(fs.readFileSync("aws.json"));

aws.config.update({accessKeyId: awsConfig.awsKey, secretAccessKey: awsConfig.awsSecret});
aws.config.update({region: awsConfig.awsRegion});

exports.getRunningInstancesByGroup = function(groupID, cb) {
	
	var filter = {
		Filters: [
			{
				Name: "group-name",
				Values : [groupID]
			},
			{
				Name: "instance-state-name",
				Values : ["running"]
			}
		]
	}
	
	new aws.EC2().describeInstances(filter, function(err, data) {
		if (err != null)
			cb("Error retrieving instance list:"+err, null)
			
		var reservationList = data.Reservations;
		
		var retList = [];
		
		for(i = 0; i < reservationList.length; i++) {
			
			var instanceList = reservationList[i].Instances;
			
			for(j = 0; j < instanceList.length; j++) {
				var item = {};
				item.instanceID = instanceList[j].InstanceId;
				item.imageID = instanceList[j].ImageId;
				item.state = instanceList[j].State;
				item.privateIpAddress = instanceList[j].PrivateIpAddress;
				item.publicIpAddress = instanceList[j].PublicIpAddress;
				item.securityGroupID = instanceList[j].SecurityGroups[0].GroupId;
				item.securityGroupName = instanceList[j].SecurityGroups[0].GroupName;
				item.keyName = instanceList[j].KeyName;
				item.launchTime = instanceList[j].LaunchTime;
				retList.push(item);
			}
		}
		
		cb(null, retList);
		
	});

}

exports.isInstanceRunning = function(instanceID, cb) {
	var filter = {
		Filters: [
			{
				Name: "instance-id",
				Values : [instanceID]
			}
		]
	}
	
	new aws.EC2().describeInstances(filter, function(err, data) {
		if (err != null)
			cb("Error retrieving instance details:"+err, null)
			
		var instance = data.Reservations[0].Instances[0];
		
		var status = {
			running: instance.State == "running"?true:false,
			publicIpAddress: instance.PublicIpAddress,
			privateIpAddress: instance.PrivateIpAddress
		}
		
		if(instance.State == "running")
			cb(null, {running:true})
		else
			cb(null, {running:false})
		
	});
}

exports.terminateInstance = function(instanceID, cb) {
	
	new aws.EC2().terminateInstances({InstanceIds : [instanceId]}, function(err, data) {
		if (err != null)
			cb("Error terminating instances:",instanceIdList, ":", err)
			
		cb();
		
	});

}

exports.runInstances = function(count, userDataList, cb) {

	var config = {
		ImageId : awsConfig.ec2ImageID,
		MinCount : count,
		MaxCount: count,
		KeyName : awsConfig.ec2KeyName,
		SecurityGroups  : [awsConfig.ec2GroupName],
		InstanceType : awsConfig.ec2InstanceType,
		UserData : new Buffer(userDataList.join("\n").toString('base64'))
	}
	
	new aws.EC2().runInstances({InstanceIds : [instanceId]}, function(err, data) {
		if (err != null)
			cb("Error launching instances:",instanceIdList, ":", err)
		
		var instanceIDList = [];
		for(i = 0; i < data.Instances.length; i++) {
			instanceIDList.push(data.Instances[i].InstanceId)
		}
		cb(null, instanceIDList);
	});

}