var	fs = require("fs");
var	crypto = require("crypto");
var JSZip = require('jszip');
var S = require('string');
var os = require('os');
var path = require('path');

module.exports = function() {

	// constructor
	var _fileList = [];
	
	// private methods
	var getFilesInFolder = function(srcFolder) {
		
		var folderEntryList = fs.readdirSync(srcFolder);
		
		var retList = [];
		
		for(var i = 0; i < folderEntryList.length; i++) {
			var folderEntry = srcFolder + "/" + folderEntryList[i];
			
			if(folderEntryList[i].substring(0, 1) === ".") {
				continue;
			}
			
			var stat = fs.statSync(folderEntry);
			
			if(stat.isFile()) {
				retList.push(folderEntry);
			} else  {
				retList = retList.concat(getFilesInFolder(folderEntry));
			}
		}
		
		return retList;
	}
	
	
	return {
		// public methods
		
		addFile: function(srcName, destName) {
			if(!fs.existsSync(srcName) || !fs.statSync(srcName).isFile())
				throw new Error(srcName + " does not exist or it's not a file");
				
			_fileList.push({src:srcName, dest:destName})
		},
		
		addFolder: function(srcFolder, destFolder) {
			if(!fs.existsSync(srcFolder) || !fs.statSync(srcFolder).isDirectory())
				throw new Error(srcName + " does not exist or is not a folder");
				
			var fileList = getFilesInFolder(srcFolder);
			
			for(var i = 0; i < fileList.length; i++) {
				var singleFile = fileList[i];
				
				var singleFileDestName = destFolder + S(singleFile).chompLeft(srcFolder).s;
				
				_fileList.push({src: singleFile, dest: singleFileDestName});
				
			}
			
		},
		
		toTempFile: function() {
			
			var zip = new JSZip();
			
			for(var i = 0; i < _fileList.length; i++) {
			
				zip.file(_fileList[i].dest, fs.readFileSync(_fileList[i].src));
			
				//zip.addFile(_fileList[i].dest, fs.readFileSync(_fileList[i].src), "");
			}
			
			var ret = os.tmpdir() + path.sep + crypto.randomBytes(10).readUInt32LE(0) + '.zip'
			
			var content = zip.generate({compression: "DEFLATE", type:"nodebuffer"});
			
			fs.writeFileSync(ret, content);
			
			return ret;
		}
		
		
	}

}