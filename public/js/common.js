$(function () {
    // $('.footer').load('./common/footer.html')
    // $('.header').load('./common/header.html')
  })

saveLocalStorage = ( name, jsstr ) => {
    localStorage.setItem(name, JSON.stringify(jsstr));
}

getLocalStorage = ( name ) => {
    return  JSON.parse( localStorage.getItem(name) );
}

  
  function callAjax(fucname, data, id) {
    var res;

    let urls = window.location.protocol + "//" + window.location.hostname + ":9022/" + fucname

    $.ajax({
      url:   urls,
      method: "POST",
      data: JSON.stringify(data),
      dataType: "json",
      success: function (response) {
        saveLocalStorage('hpcstorage', response)
        if (response[0]['Id'] == data['Id'])
            location.href = './jobmanager.html?server=' + response[0]['Server'];
        else
          alert('login falure')
      },
      error: function (jqXHR, textStatus, errorThrown) {
        if (textStatus == "error") {
          return "error has occurred: " + jqXHR.status + " " + jqXHR.statusText
        }
      }
    });
  
    return res;
  }