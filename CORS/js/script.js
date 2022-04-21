function play() {
	let mode = document.getElementById("mode").value;
	let order = document.getElementById("order").value;
	send(mode, order)
};
	
async function send(mode, order) {
	let res = await fetch("http://localhost:8080/", {
		method: 'POST',
		mode: 'cors',
		cache: 'no-cache',
		credentials: 'same-origin',
		headers: {'Content-Type': 'application/json'},
		redirect: 'follow',
		referrerPolicy: 'no-referrer',
		body: '{"mode":"'+mode+'","order":"'+order+'"}'});
	let inf = await res;
	inf.text().then(function(text) {
		text=JSON.parse(text);
		document.write(text.html);
		if (text.op!="") {
			document.getElementById(text.op).style.display="none";
		}
	});
};
