$(document).ready(function(){
    console.log("This is on start up")
    $("button").click(function(){
      var myForm = document.createElement("FORM");
      myForm.setAttribute("action","/gameSelected");
      myForm.setAttribute("method", "POST");
      
      var input = document.createElement("INPUT");
      input.setAttribute("type", "text");
      input.setAttribute("value", $(this).val());
      myForm.appendChild(input)
      $(document.body).append(myForm);
      console.log("This is the value ", $(this).val());
      $(myForm).submit();
    });
  });