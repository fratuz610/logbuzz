'use strict';

var services = angular.module('logbuzz.services', []);


services.factory('logService', ['$timeout', '$http', '$rootScope', 'logbuzzURL',
    function($timeout, $http, $rootScope, logbuzzURL) {
        var ret = {

            latestTimestamp: 0,

            updateLogList: function(skip, level, tagList) {

                var params = {};
                if(skip != undefined) params.skip = 500*skip;
                if(tagList != undefined && tagList != "") params.tagList = tagList;
                if(level != undefined && level != "") params.level = level;

                $http.get(logbuzzURL+"/api/query", {params:params})
                    .success(function(data, status, headers, config) {
                        console.log("logItems: " + data.length);
                        $rootScope.$emit('newLogItemListAvailable', data);
                    });
            }

        };

        return ret;
    }]);