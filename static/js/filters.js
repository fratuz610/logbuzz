'use strict';

var filters = angular.module('logbuzz.filters', []);

filters.filter('dateTime', function() {
	return function(dateTime, param) {
		if(param == "includeTimezone")
			return moment(dateTime).format("YYYY-MM-DD HH:mm:ss Z")
		else
			return moment(dateTime).format("YYYY-MM-DD HH:mm:ss")
	};
});

filters.filter('timeSpan', function() {
	return function(timeSpan) {
		return moment().subtract('ms', timeSpan).fromNow();
	};
});

filters.filter('calendar', function() {
    return function(ts) {
        return ts>0?moment(ts).format('lll'):"";
    };
});

filters.filter('available', function() {
    return function(data) {
        return (data == null || data == undefined || data == "")?"N/A":data;
    };
});

filters.filter('toLowerCase', function() {
	return function(str) {
		return str.toLowerCase();
	};
});

filters.filter('toUpperCase', function() {
    return function(str) {
        return str.toUpperCase();
    };
});

filters.filter('join', function() {
	return function(list, separator) {
		return list.join(separator);
	};
});

filters.filter('range', function() {
  return function(input, total) {
    total = parseInt(total);
    for (var i=0; i<total; i++)
      input.push(i);
    return input;
  };
});