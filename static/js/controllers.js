var app = angular.module('logbuzz', ['logbuzz.filters', 'logbuzz.services', 'logbuzz.directives']);

app.config(['$routeProvider', function($routeProvider) {
    $routeProvider.
        when('/log', {templateUrl: 'partials/log.html', controller: 'LogCtrl'}).
        otherwise({redirectTo: '/log'});
}]);

app.controller('LogCtrl', ['$scope','$rootScope', '$location', 'logService',
    function($scope, $rootScope, $location, logService) {

        $scope.loaded = false;
        $scope.logItemList = [];

        $scope.skip = 0;
        $scope.level = "";
        $scope.tagList = "";

        $rootScope.$on('newLogItemListAvailable', function(e, data) {
          $scope.logItemList = data;
          $scope.loaded = true;
        });

        $scope.skipBack = function() {
            if($scope.skip <= 0)
                return;

            $scope.skip--;
            $scope.refresh();
        };

        $scope.skipForward = function() {
            if($scope.logItemList.length == 0)
                return;

            $scope.skip++;
            $scope.refresh();
        };

        $scope.setLevel = function(newLevel) {
            $scope.level = newLevel;
            $scope.refresh();
        }

        $scope.refresh = function() {
            $scope.loaded = false;
            logService.updateLogList($scope.skip, $scope.level, $scope.tagList)
        }

        $scope.levelString = function() {
            if($scope.level == null || $scope.level == "")
                return "all";
            return $scope.level.toUpperCase();
        }
	}
]);


app.run(function($rootScope, $location, $routeParams, logService) {
	console.log("Application bootstrapping"); 


	$rootScope.$on('$routeChangeSuccess', function(next, current) {
		
        console.log("$routeChangeSuccess: current: %o", current);

        logService.updateLogList();
	});


	
});

app.factory('logbuzzURL', function($location) {
    //return 'http://127.0.0.1:1212';
    return "";
});
//app.value('logbuzzURL', 'http://127.0.0.1:1212');