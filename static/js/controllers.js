var app = angular.module('logbuzz', ['logbuzz.filters', 'logbuzz.services', 'logbuzz.directives']);

app.config(['$routeProvider', function($routeProvider) {
    $routeProvider.
        when('/log', {templateUrl: 'partials/log.html', controller: 'LogCtrl'}).
        otherwise({redirectTo: '/log'});
}]);

app.controller('LogCtrl', ['$scope','$rootScope', '$location', 'logService',
    function($scope, $rootScope, $location, logService) {

        /*
        $scope.savedLogins = loginService.loadSavedLogins();

        $scope.loginUsername = "";
        $scope.loginPassword = "";
        $scope.loginToken = "";
        $scope.isLoading = false;
        $scope.validFor = 60*60*24;

        $scope.loginCred = function() {
            $scope.loginError = "";
            $scope.isLoading = true;
            loginService.loginCred($scope.loginUsername, $scope.loginPassword, $scope.validFor);
        }

        $scope.loginToken = function(tokenID) {
            $scope.loginError = "";
            $scope.isLoading = true;
            console.log("Calling with tokenID: %o", tokenID);
            loginService.loginToken(tokenID);
        }

        $rootScope.$on('loginCredSuccess', function(e, data) {
          loginService.saveLogin(data.email, data.tokenID);
          $rootScope.email = data.email;
          $rootScope.tokenID = data.tokenID;
          $scope.isLoading = false;
          $scope.$emit('closeModal',function() {
            $location.path("/units");
            $scope.$apply();
          });
        });

        $rootScope.$on('loginCredFailure', function(e, errorCode, message) {
            $scope.loginError = "Could not log you in: " + errorCode + " - " + message;
            $scope.isLoading = false;
        });

        $rootScope.$on('loginTokenSuccess', function(e, data) {
          $rootScope.email = data.email;
          $rootScope.tokenID = data.tokenID;
          $scope.isLoading = false;
          $scope.$emit('closeModal',function() {
            $location.path("/units");
            $scope.$apply();
          });
        });

        $rootScope.$on('loginTokenFailure', function(e, errorCode, message, tokenID) {
            $scope.loginError = "Invalid/Expired token: " + errorCode + " - " + message;
            console.log("Token login failed: %o", tokenID);
            loginService.deleteLogin(tokenID);
            $scope.isLoading = false;
        });
         */
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