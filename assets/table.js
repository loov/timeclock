(function(Table) {
	var MIN_HOURS = 0;
	var MAX_HOURS = 20;

	Table.ClearRow = function(event) {
		if (!event) return;
		event.preventDefault();
		event.stopPropagation();

		var target = event.target;
		if (!target) return;

		var row = target;
		while (row) {
			if (row.tagName === "TR") break;
			row = row.parentElement;
		}
		if (!row) return;

		var fields = row.querySelectorAll("input,select");
		for (var i = 0; i < fields.length; i++) {
			var field = fields[i];
			field.value = "";
		}
	};

	Table.RoundUpHours = function(event) {
		if (!event) return;
		event.preventDefault();
		event.stopPropagation();

		var target = event.target;
		if (!target) return;
		if (target.tagName != "INPUT") return;

		var hours = parseFloat(target.value);
		if (isNaN(hours) || (hours <= MIN_HOURS) || (hours >= MAX_HOURS)) {
			target.value = "";
			return;
		}

		var hours2 = Math.ceil(hours * 2);
		if (hours2 % 2 == 0) {
			target.value = hours2 / 2;
		} else {
			target.value = (hours2 / 2).toFixed(1);
		}
	};
})(window.Table = {});