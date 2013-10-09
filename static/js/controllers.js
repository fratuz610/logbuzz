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

        $rootScope.$on('newLogItemListAvailable', function(e, data) {
          $scope.logItemList = data;
          $scope.loaded = true;
        });
	}
]);


app.run(function($rootScope, $location, $routeParams, logService) {
	console.log("Application bootstrapping"); 
	
	$rootScope.$on('$routeChangeSuccess', function(next, current) {
		
        console.log("$routeChangeSuccess: current: %o", current);

        logService.start();
	});
	
});

app.factory('logbuzzURL', function($location) {
    return 'http://127.0.0.1:1212';
});
//app.value('logbuzzURL', 'http://127.0.0.1:1212');