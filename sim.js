// Original script for the Yard Sale Model, h/t https://physics.umd.edu/hep/drew/math_general/yard_sale.html
var gain;
var loss;
var initial_amount;
var plays;
var people;
var ymax;
var niter;
var opts, opts2;
var xt = [];
var xit = [];
var it;
var nit;
var x_in = [];		// input to Fisher-Yates
var x2 = [];		// temporary for squeezing out zeros
var v =[];			// how much each person has

var ystimer = null;
var ystime = 200;
function run_anim() {
	plays = Number(document.getElementById("plays").value);
	people = Number(document.getElementById("people").value);
	gain = Number(document.getElementById("gain").value);
	loss = Number(document.getElementById("loss").value);
	ymax = people*initial_amount;
	opts_all.ymax = ymax;
//	opts_one.ymax = initial_amount;
	x_in.length = 0;
	x2.length = 0;
	v.length = 0;
	for (var k=0; k<people; k++) {
		x_in.push(k);
		x2.push(0);
		v.push(initial_amount);
	}
    if (ystimer == null) ystimer = setInterval("ys_iterate()",ystime);
}
function stop_anim() {
	if (ystimer == null) return;
	else clearInterval(ystimer);
	ystimer = null;
}
function ys_iterate() {
	yard_sale_iteration();
	plotit();
	niter = niter + plays;
	document.getElementById("nanim").innerHTML = niter;
}

function yard_sale_iteration() {
	//
	// x_in is an array that contains an index.   v[x_in[i]] is the value for person x_in[i]
	//
//	console.log("1: "+x_in);
//	console.log("   "+v);
	var debug = false;
	for (var j=0; j<plays; j++) {
		var debug = false;
		var halfpeople = people/2;
		if (people & 1 == 1) {
			//
			// oops, it's odd.   so we will just have to not pair off the
			// last one and let it slide for this play
			//
			halfpeople = (people-1)/2;
		}
//		console.log("j="+j+" people="+people+" half="+halfpeople);
		//
		// first reorder
		//
//		console.log(x_in);
		fisher_yates();
		//
		// now use the shuffled x_in and loop over pairs and play the game
		//
		xt[nit] = nit;
		xit[nit] = v[it];
		nit++;
		for (var i=0; i<halfpeople; i++) {
			let i1 = x_in[i];
			let i2 = x_in[i+halfpeople];
//			console.log("i="+i+"  x_in["+i+"]="+i1+"  v["+i1+"]="+v[i1]);
			let xr = Math.random();
			let win = xr < 0.5;
			let v1 = v[i1];
			let v2 = v[i2];
			//
			// j1 points to the wealthier of the 2, j2 to the poorer one
			//
			let j1 = i1;
			let j2 = i2;
			if (v2 > v1) {
				j1 = i2;
				j2 = i1;
			}
			if (win) {
				//
				// a "win" means the poorer person gains "gain" % of their wealth,
				// coming from the hide of the richer person.
				//
				let delta = v[j2]*gain/100;
				v[j2] = v[j2] + delta;
				v[j1] = v[j1] - delta;
			}
			else {
				//
				// a "loss" means the poorer person loses "loss" % of their wealth,
				// transfering it to the ricer person
				//
				let delta = v[j2]*loss/100;
				v[j1] = v[j1] + delta;
				v[j2] = v[j2] - delta;
			}
//			console.log("  i="+i+" i1/i2="+i1+"/"+i2+" before: v[i1]/v[i2]="+v1+"/"+v2+
//				" after: v[i1]/v[i2]="+v[i1]+"/"+v[i2]);
		}
		//
		// calculate the sum (should be constant)
		//
		let sum = 0;
		for (var k=0; k<people; k++) {
			let ki = x_in[k];
			sum = sum + v[ki];
		}
		document.getElementById("spansum").innerHTML = Math.floor(sum);
//		console.log(" Sum = "+sum);
		//
		// all done for this play
		//
	}
//	console.log(v)
}

function plotit() {
	//
	// plot the distribution of wealth
	//
	plot_wealth();
	//
	// plot the distribution for a single person to see how it changes
	//
	plot_single();
}

function plot_wealth() {
//	console.log("x2.length="+x2.length+" v.length="+v.length);
	//
	// now just plot things sequencially.   
	//
	x2.length = 0;
	for (var i=0; i<people; i++) x2.push(i);
	umd_scatterplot("canvas","Yard Sale Model",x2,v,opts_all);
}

function plot_single() {

//	console.log("x2.length="+x2.length+" v.length="+v.length);
	//
	// now just plot things sequencially.   
	//
	umd_scatterplot("single","Sequence for person "+it.toString(),xt,xit,opts_one);
}

function fisher_yates() {
	//
	// takes a set of numbers in xn and reorders in a random way
	//
	len = x_in.length;
//	console.log("FY: len="+len);
	var xtemp;
	for (var i=len-1; i>0; i--) {
		//
		// generate a random integer between 0 and i inclusive [0,i]
		//
//		console.log("i="+i);
		let r = Math.random();
		let j = Math.floor(r*(i+1));
//		console.log(" xn["+i+"]="+x_in[i]+"  r="+r+"  xn["+j+"]="+x_in[j]);
		xtemp = x_in[i];
		x_in[i] = x_in[j];
		x_in[j] = xtemp;
	}
}

function reset_ys() {
	init_ys();
}

function init_ys() {
	people = 100;
	plays = 100;
	initial_amount = 100;
	gain = 20;
	loss = 17;
	niter = 0;
	it = Math.floor(Math.random()*people);
	document.getElementById("plays").value = plays;
	document.getElementById("people").value = people;
	document.getElementById("spansum").innerHTML = people*initial_amount;
	document.getElementById("gain").value = gain;
	document.getElementById("loss").value = loss;
	ymax = people*initial_amount;
	x_in.length = 0;
	x2.length = 0;
	v.length = 0;
	xt.length = 0;
	xit.length = 0;
	nit = 0;
	for (var k=0; k<people; k++) {
		x_in.push(k);
		x2.push(0);
		v.push(initial_amount);
	}
	//
	// options for the 2 plots:
	//
	opts_all = {
		xtitle :'n',
		width: 500,
		height: 500,
		ymin: 0.0,
		xmin: 0,
		ymax: ymax,
		ygrid: 5,
		xmax: people,
		active : true,
		active_div: "active",
		background: "#b5b5b5",
		show_stats: false,
		plot_type: "bar",
		bar_type: "open",
		line_color: "blue"
	}
	opts_one = {
		xtitle :'Sequence Number',
		width: 500,
		height: 500,
		ymin: 0.0,
		xmin: 0,
		ygrid: 5,
		active : true,
		active_div: "single",
		background: "#b5b5b5",
		show_stats: false,
		plot_type: "bar",
		bar_type: "open",
		line_color: "blue"
	}	
	plotit();
}
init_ys();

