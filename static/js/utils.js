function htmlEncode(str) {

    if(typeof(str) == "undefined" || str == null) {
        return "";
    }

    if(typeof(str) != "string") {
        str = str.toString();
    }

    return str.replace(/&/g, '&amp;')
        .replace(/'/g, '&#39;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;');
}

function ajaxRequest(method,url,params,callback){
    jQuery.ajax({
        type: method,
        url: url,
        data: params,
        beforeSend: function(xhr){
            //使用jQury.base64 传xsrf值
            xhr.setRequestHeader("X-Xsrftoken",jQuery.base64.atob(jQuery.cookie("_xsrf").split("|")[0]))
        },
        success: function (response){
            switch (response["code"]) {
                case 200:
                    callback(response);
                    alert(response["text"])
                    break;
                case 400:
                    var err = [];
                    jQuery.each(response["result"],function (k,v){
                        err.push(v['Message']);
                    });
                    if(!err) {
                        err.push(response['text']);
                    }
                    alert(err.join("\n"))
                    break;
                case 401:
                    alert(response['text'])
                    window.location.replace("/")
                    break;
                case 500:
                    alert(response['text']);
                    break;
                default:
                    alert(response['text']);
                    break;
            }
        },
        dataType: "json"
    });
}