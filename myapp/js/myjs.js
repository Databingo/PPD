        var ws;
        var record = []; 
        var count = 0; 
        var last_show_value = "";
        var last_input_value = "";

	function start_websocket(){

            ws = new WebSocket('ws://localhost:8300/search');

            //---- 
            ws.onopen = function(evt){
                console.log('open');
                setInterval(watch_input, 600);
            };
            
            //---- 
            ws.onmessage = function(evt){
                var pre = document.getElementById("screen_1").innerHTML;
                console.log(evt.data);
	        try {var msg = JSON.parse(evt.data)} catch (e){
                    console.log(e);
	        };
                var input_val = document.getElementById("search").value;
	        if (msg && msg.word == input_val && msg.type == 0){
	            console.log("------")
	            console.log(evt.data.length <= last_show_value.length)
		    console.log((last_show_value == "" || (evt.data.length <= last_show_value.length)))
	            console.log(evt.data.length, last_show_value.length)
		    count ++
	            console.log(input_val, count)
	            console.log("------")
	            if (last_show_value == "" || (evt.data.length <= last_show_value.length)) {
                        //document.getElementById("screen_1").innerHTML=msg[msg.lg]+"<hr class='hr'>"+msg[(msg.lg=="pl")?"cn":"pl"];
                        document.getElementById("screen_1").innerHTML=msg[msg.lg]+"<hr class='hr'>";
                        document.getElementById("screen_2").innerHTML=msg[(msg.lg=="pl")?"cn":"pl"]+"<hr class='hr'>";
		        // console.log(msg)
                        last_show_value  = evt.data;
	            }
	            record.push(msg);
	        }


	        if (msg && msg.word == input_val && msg.type == 1){
                        document.getElementById("hit_ranking").innerHTML= document.getElementById("hit_ranking").innerHTML + "<hr class='hr' >"+ msg.hit+" "+msg.count;
		}
            };

            //---- 
            ws.onclose= function(evt){
            console.log('close');
	        ws = null;
	        setTimeout(start_websocket, 2000)
            };
        };
        start_websocket()

        function watch_input(){
            var input_val = document.getElementById("search").value;
	    if (input_val != last_input_value){
	        if (input_val){
	            ws.send(input_val);
	            }
                // clear
                // document.getElementById("screen_1").innerHTML = "";
                document.getElementById("hit_ranking").innerHTML = "";
                last_show_value = "";
                record = [];
	        count = 0;
                last_input_value = input_val;
	    };
	};


       function input_pali_alphabet(alph){
            var input_val = document.getElementById("search").value;
	    input_val = input_val + alph
            document.getElementById("search").value=input_val;
            document.getElementById("search").focus();
       };
