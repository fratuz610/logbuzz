'use strict';

var services = angular.module('logbuzz.services', []);


services.factory('logService', ['$timeout', '$http', '$rootScope', 'logbuzzURL',
    function($timeout, $http, $rootScope, logbuzzURL) {
        var ret = {

            self:this,

            latestTimestamp: 0,

            _cancelRefresh: null,

            start: function() {
                console.log("Start called: " + self._cancelRefresh);

                self.latestTimestamp = 0;

                self._cancelRefresh = $timeout(function myFunction() {
                    $http.get(logbuzzURL+"/api/stats")
                        .success(function(data, status, headers, config) {

                            if(data.latestTimestamp > self.latestTimestamp) {

                                console.log("Updating timestamp " + data.latestTimestamp);
                                self.latestTimestamp = data.latestTimestamp;

                                $http.get(logbuzzURL+"/api/query")
                                    .success(function(data, status, headers, config) {
                                        console.log("logItems: " + data.length);
                                        $rootScope.$emit('newLogItemListAvailable', data);
                                    });
                            }
                            self._cancelRefresh = $timeout(myFunction, 1000);
                        });

                },1000);

            },
            stop: function() {
                console.log("Stop called: " + self._cancelRefresh);
                $timeout.cancel(self._cancelRefresh);
                self._cancelRefresh = null;
            },

            isRunning: function() {
                return self._cancelRefresh != null;
            }
        };

        return ret;
    }]);