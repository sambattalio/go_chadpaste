
function button_upload(ev) {
    let file = ev.target.files[0];
    let formData = new FormData();

    formData.append("file", file);
    formData.append("expiration", 3000);
    console.log("here");
    fetch('/post', {method: "POST", body: formData})
      .then(function(response) {
        return response.text();
    }).then(function(data) {
        window.location.href = '?' + data;
    });
    console.log("here 2");
}


function get_query() {
  var query = window.location.search.substring(1);

  // If there was some kind of query string
  if (query) {
    var file_location = 'f/' + query;

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
};

window.addEventListener("dragover",function(e){
    e = e || event;
    e.preventDefault();
},false);
window.addEventListener("drop",function(e){
    e = e || event;
    e.preventDefault();
},false);

function dragover_handler(ev) {
    ev.preventDefault();
}
window.doDrop = function(ev) {
    var dti = ev.dataTransfer.files;
    if (dti === undefined) {
        console.log("DataTransferItem NOT supported.");
        console.log("DataTransferItemList NOT supported.");
    } else {
    }
}
