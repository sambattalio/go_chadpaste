
function button_upload() {
    let file = document.getElementById('file_drop').files[0];
    console.log(file)
    let test = new File([file], file.name, { type: file.type });
    let formData = new FormData();

    formData.append("file", test);
    formData.append("expiration", document.getElementById('expiration').value);

    var radios = document.getElementsByName('expir');
    for (var i = 0, length = radios.length; i < length; i++) {
      if (radios[i].checked) {
        // do whatever you want with the checked radio
        formData.append("expiration_type", radios[i].value);

        // only one radio can be logically checked, don't check the rest
        break;
      }
    }
    fetch('/post', {method: "POST", body: formData})
      .then(function(response) {
        return response.text();
    }).then(function(data) {
       window.location.href = '?' + data;
    }); 
}

function get_query() {
  var query = window.location.search.substring(1);
  // If there was some kind of query string
  if (query) {
    var file_location = 'f/' + query;

    // delete expiration option and replace with info
    document.getElementById('submitForm').remove();

    var expir_info = document.createElement("text");
    document.body.appendChild(expir_info);
    fetch("/expir/" + query, {method: "GET"})
      .then(function(response) {
        return response.json();
      }).then(function(data) {
        if (data['type'] == -2) {
            window.location.href = '/';
        } else if (data['type'] == 1) {
            expir_info.innerText = "Expires in " + data['value'] + " views";
        } else if (data['type'] == 0) {
            var date = new Date(data['value'] * 1000);
            console.log(date);
            var now  = new Date();
            date = date.getTime() - now.getTime()
            if (date < 0) {
                window.location.href = '/';
            }
            time = date / 1000;

            expir_info.innerText = "Expires in: " + time + "seconds";
        } else {
        }
      });

    var raw_btn = document.createElement("BUTTON");
    raw_btn.innerText = "RAW FILE";
    raw_btn.onclick = function () {
        window.location.href = '/f/' + query;
    }
    raw_btn.style.width = 100;
    raw_btn.style.height = 20;
    raw_btn.padding = 30;
    raw_btn.paddingRight = 300;
    raw_btn.style.maxWidth = "25%";
    raw_btn.style.maxHeight = "5%";
    // document.body.insertBefore(raw_btn, document.getElementById("content"));
    document.body.appendChild(raw_btn);
    document.body.appendChild(document.createElement("BR"));
    // Handle Images
    if (query.endsWith("jpg") || query.endsWith("jpeg") || query.endsWith("png") || query.endsWith("gif")) {
      var img   = new Image();
      if (query) {
        img.src   = file_location;
        img.id    = 'uploaded';
	img.style.maxWidth = '100%';
        document.body.appendChild(img);
        //document.getElementById("content").appendChild(img);
        return true;
      }
    } else if (query.endsWith("mp4") || query.endsWith("webm")) {
        var vid = document.createElement("VIDEO");
        vid.src = file_location;
        vid.id = 'uploaded';
        vid.controls = true;
        vid.style.maxWidth = '100%';
        document.body.appendChild(vid);
        return true;
    } else {
      // Assume it is text of some kind
      // TODO add raw button
      fetch(file_location)
        .then(response => response.text())
        .then(text => {
          var pre = document.createElement("pre");
          var code = document.createElement("code");
          var newDiv = document.createElement("div");
          newDiv.style = "white-space: pre-wrap";
          var newContent = document.createTextNode(text.substring(0, 42069));
          code.appendChild(newContent);
          pre.appendChild(code);
          newDiv.appendChild(pre);
          document.body.appendChild(newDiv);
        })
      return true;
    }
  }
	return false;
}


window.onload = function() {
    if (get_query()) {
        $('#box').remove();
        $('#button').remove();
    }
}
