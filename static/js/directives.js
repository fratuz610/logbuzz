'use strict';

var directives = angular.module('logbuzz.directives', []);

directives.directive('popup', function factory($location) {
    return {
        restrict: 'A',
        link: function (scope, element, attrs) {
          console.log("Calling modal with attrs %o", attrs);

          $(element).modal({
              backdrop: "static",
              keyboard: false
          });

          scope.$on('closeModal', function(e, fn) {
            $(element).modal('hide');
            $(element).on('hidden.bs.modal', function () { fn(); })

          });



        }
    };
});

directives.directive('restrict', function($parse) {
    return {
        restrict: 'A',
        require: 'ngModel',
        link: function(scope, iElement, iAttrs, controller) {
            scope.$watch(iAttrs.ngModel, function(value) {
                if (!value) {
                    return;
                }
                $parse(iAttrs.ngModel)
					.assign(scope, value.toLowerCase().replace(new RegExp(iAttrs.restrict, 'g'), '').replace(/\s+/g, '-'));
            });
        }
    }
});

