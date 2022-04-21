async function use(choice) {
	let inf = await send(choice);
	let gameend = false;
	if (inf.can) {
		document.getElementById(choice).style.display="none";
		document.getElementById("altobj").innerText=inf.val;
		res = await getbotmove();
		switch (res) {
			case "win": alert("You win!");gameend=true;break;
			case "draw": alert("It is a draw!");gameend=true;break;
			case "lose": alert("You lose!");gameend=true;break;
			default:
				document.getElementById("altobj").innerText=res.val;
				document.getElementById(res.op).style.display="none";
		}
		if (gameend) {
			location.reload();
		}
	}
};
	
async function send(c) {
	let res = await fetch("http://localhost:8080/game/", {
		method: 'POST',
		mode: 'cors',
		cache: 'no-cache',
		credentials: 'same-origin',
		headers: {'Content-Type': 'application/json'},
		redirect: 'follow',
		referrerPolicy: 'no-referrer',
		body: c});
	let inf = await res;
	return inf.text().then(function(text) {
		return JSON.parse(text);
	});
};

async function getbotmove() {
	let res = await fetch("http://localhost:8080/game/", {
		method: 'GET',
		mode: 'cors',
		cache: 'no-cache',
		credentials: 'same-origin',
		headers: {'Content-Type': 'application/json'},
		redirect: 'follow',
		referrerPolicy: 'no-referrer'});
	let inf = await res;
	return inf.text().then(function(text) {
		switch (text) {
			case "win":case "draw":case "lose": return text; break;
			default: return JSON.parse(text);
		}
	});
}
