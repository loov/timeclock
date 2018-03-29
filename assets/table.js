(function(Table) {
	Table.ClearRow = function(event) {
		if (!event) return;

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
})(window.Table = {});