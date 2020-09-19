$(document).ready(function(){
    console.log("This is on start up")
    const url = "127.0.0.1:8080"

    $("#return").click(function(){
      var myForm = document.createElement("FORM");
      myForm.setAttribute("action","/gameSelected");
      myForm.setAttribute("method", "POST");
      
      var input = document.createElement("INPUT");
      input.setAttribute("type", "text");
      input.setAttribute("value", $(this).val());
      myForm.appendChild(input)
      $(document.body).append(myForm);
      $(myForm).submit();
    });

    // $(".light, .dark").click(function(){
    //   var myForm = document.createElement("FORM");
    //   myForm.setAttribute("action","/gameSelected");
    //   myForm.setAttribute("method", "POST");
      
    //   var input = document.createElement("INPUT");
    //   input.setAttribute("type", "text");
    //   input.setAttribute("name", $(this).prop("name"))
    //   input.setAttribute("value", $(this).prop("name") + " " + $(this).val() + " piece clicked");
    //   myForm.appendChild(input)
    //   $(document.body).append(myForm);
    //   $(myForm).submit();
    // });

    $(".light, .dark").click(function(){
      var myForm = document.createElement("FORM");
      // myForm.setAttribute("action","/gameSelected");
      myForm.setAttribute("method", "POST");
      
      var input = document.createElement("INPUT");
      input.setAttribute("type", "text");
      input.setAttribute("name", $(this).prop("name"))
      input.setAttribute("value", $(this).prop("name") + " " + $(this).val() + " piece clicked");
      myForm.appendChild(input)
      $(document.body).append(myForm);
      // $(myForm).submit();
      fetch("/game", {
          method:"POST",
          body: new FormData(myForm)
      }).then(
        response => response.text()
      ).then(
        (data) => {console.log("Reponse from server: ", data)}
      ).catch(
        error => console.error("Error encountered: ", error)
      )

    });


    $("button").click(function () {
      console.log("Value is: ", $(this).val())
      console.log("Name is: ", $(this).prop("name"))
      console.log("Class is: ", $(this).prop("class"))      
      console.log("Button Text is: ", $(this).text(), $(this).text().length)
    })


    
  });