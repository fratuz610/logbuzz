var async = require("async"),
    path = require("path"),
	aws = require("aws-sdk"),
    fs = require("fs");

var config = JSON.parse(fs.readFileSync("config.json"));
var awsConfig = JSON.parse(fs.readFileSync("aws.json"));

aws.config.update({accessKeyId: awsConfig.awsKey, secretAccessKey: awsConfig.awsSecret});
aws.config.update({region: 'ap-southeast-2'});
	
var uploadToS3Testing = function() {

	var options = {}
	options.Bucket = awsConfig.s3Bucket;
	options.Key = awsConfig.s3Key;
	
	var s3 = new aws.S3();
	
	// first we delete the old file
	
	s3.deleteObject(options, function(err, data) {
		if(err)
			console.log("Unable to delete ", options.Bucket, "/", options.Key, " because: ", err);
		else
			console.log("Previous version of ", options.Bucket, "/", options.Key, " removed");
	});
	
	if(!fs.existsSync(config.logbuzzExecutable))
		throw new Exception('File ', config.logbuzzExecutable, " does not exist, unable to continue")
	
	options.ACL = "public-read";	
	options.Body = fs.readFileSync(config.logbuzzExecutable);
	
	s3.putObject(options, function(err, data) {
	  
		if (err != null)
			throw new Exception('Error uploading file ', config.logbuzzExecutable, " because: ", err)
		
		console.log("Upload completed");
	});
	
	console.log("Upload scheduled");
};

uploadToS3Testing();