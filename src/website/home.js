$(document).ready(function () {
  console.log("STARTED");
  var localText = "Start Game";
  var aiText = "Battle the AI";
  var onlineText = "Find an Opponent";
  var inner = ""
  var clicked = false
  document.getElementById("startGame").innerHTML = onlineText

  $("#local").click(function () {
    document.getElementById("startGame").innerHTML = localText;
  });

  $("#ai").click(function () {
    document.getElementById("startGame").innerHTML = aiText;
  });

  $("#online").click(function () {
    if (clicked) {
    document.getElementById("startGame").innerHTML =
    "<span class='spinner-border spinner-border-sm mr-2' role='status' aria-hidden='true'></span>" +
    onlineText;
    } else {
        document.getElementById("startGame").innerHTML = onlineText;
        
    }
  });

  $("#startGame").click(function () {
    inner = document.getElementById("startGame").innerHTML;
    if (inner == onlineText) {
        clicked = true
      document.getElementById("startGame").innerHTML =
        "<span class='spinner-border spinner-border-sm mr-2' role='status' aria-hidden='true'></span>" +
        inner;
    }
  });

 $("form").submit(function (e) {
    console.log("REACHED ")

    e.preventDefault();
    var option = 0;
    console.log("text", document.getElementById("startGame").innerHTML)
    switch (document.getElementById("startGame").innerHTML) {
      case localText:
        option = 0;
        break
      case aiText:
        option = 1;
        break
      default:
        option = 2;
    }
    var myForm = document.createElement("FORM");
    myForm.setAttribute("method", "POST");
    var input = document.createElement("INPUT");
    input.setAttribute("type", "text");
    input.setAttribute("name", "option");
    input.setAttribute("value", option);
    myForm.appendChild(input);
    console.log("About to fetch")
    console.log("option", option)
    fetch("/", {
      method: "POST",
      body: new FormData(myForm),
    })
      .then((response) => response.text())
      .then((data) => {
          console.log("Server:", data);
        //   if (!changedPage) {
              document.getElementById("gameOptions").submit()})
        //       changedPage = true
        // }})
      .catch((error) => console.error("Error encountered: ", error));
  })
});
