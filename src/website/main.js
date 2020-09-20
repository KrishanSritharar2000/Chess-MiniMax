$(document).ready(function () {
  console.log("This is on start up");
  const url = "127.0.0.1:8080";
  var movesDisplayed = false;
  var moveDisplayedPiece = "00";
  var htmlChangedPiece = [];
  var htmlChangedPieceOrigText = [];
  var whiteTurn = true;

  $("#return").click(function () {
    location.reload()
    // var myForm = document.createElement("FORM");
    // myForm.setAttribute("action", "//#endregion");
    // myForm.setAttribute("method", "GET");

    // var input = document.createElement("INPUT");
    // input.setAttribute("type", "text");
    // input.setAttribute("value", $(this).val());
    // myForm.appendChild(input);
    // $(document.body).append(myForm);
    // $(myForm).submit();
  });

  function handleResponse(response, mode, clickedButton) {
    if (mode == "opt") {
      console.log("Operating mode: opt");
      console.log("Reponse from server:", response, "length:", response.length);
      for (let i = 0; i < response.length; i += 2) {
        const htmlPiece = document.getElementById(response.substring(i, i + 2));
        console.log(
          "Pushed to array ",
          response.substring(i, i + 2),
          htmlPiece.innerHTML
        );
        htmlChangedPiece.push(response.substring(i, i + 2));
        htmlChangedPieceOrigText.push(htmlPiece.innerHTML);
        htmlPiece.innerHTML = "TAKE " + htmlPiece.innerHTML;
      }

      movesDisplayed = true;
      moveDisplayedPiece = clickedButton.prop("id");

    } else if (mode == "mov") {
      console.log("Operating mode: mov");
      console.log("Reponse from server:", response, "length:", response.length);
      var result = response.split(":")[1]
      if (result === "true") {
        clearDisplayedMoves()
        const newPiece = document.getElementById(clickedButton.prop("id"))
        const oldPiece = document.getElementById(moveDisplayedPiece)
        newPiece.setAttribute("name", oldPiece.getAttribute("name"))
        newPiece.innerHTML = oldPiece.innerHTML
        oldPiece.innerHTML = " "
        oldPiece.setAttribute("name", "empty")
        whiteTurn = !whiteTurn
      } else if (result === "false") {
        console.log("THAT WAS AN INVALID MOVE")
      }
    }
  }

  function clearDisplayedMoves() {
    for (let i = 0; i < htmlChangedPiece.length; i++) {
      const htmlPiece = document.getElementById(htmlChangedPiece[i]);
      htmlPiece.innerHTML = htmlChangedPieceOrigText[i];
    }
    htmlChangedPiece = [];
    htmlChangedPieceOrigText = [];
    movesDisplayed = false;
  }

  $(".light, .dark").click(function () {
    var myForm = document.createElement("FORM");
    myForm.setAttribute("method", "POST");

    var input = document.createElement("INPUT");
    input.setAttribute("type", "text");
    input.setAttribute("name", $(this).prop("name"));


    if (movesDisplayed) {
      if (htmlChangedPiece.includes($(this).prop("id"))) {
        console.log("Sending move request!!")
        //move the piece there
        var mode = "mov"
        input.setAttribute("value", mode + " " + moveDisplayedPiece + $(this).val());

      } else {
        //Clear displayed moves
        console.log("These are the changed pieces", htmlChangedPiece);
        clearDisplayedMoves()
        if ($(this).prop("id") == moveDisplayedPiece) {
          return;
        }
      }
    } 
    if (!movesDisplayed) {
      if (
        (whiteTurn && $(this).prop("name") == "black") ||
        (!whiteTurn && $(this).prop("name") == "white")
      ) {
        console.log("NO REQUEST TO SERVER. THAT IS OPPONENTS PIECE");
        return;
      }
      if ($(this).prop("name") == "empty") {
        console.log("NO REQUEST TO SERVER. THERE IS NO PIECE THERE");
        return;
      }
      // input.setAttribute("value", $(this).prop("name") + " " + $(this).val() + " piece selected");
      var mode = "opt"
      input.setAttribute("value", mode + " " + $(this).val());
    }
    myForm.appendChild(input);

    fetch("/game", {
      method: "POST",
      body: new FormData(myForm),
    })
      .then((response) => response.text())
      .then((data) => handleResponse(data, mode, $(this)))
      .catch((error) => console.error("Error encountered: ", error));
  });

  $("button").click(function () {
    console.log("Value is: ", $(this).val());
    console.log("Name is: ", $(this).prop("name"));
    console.log("Class is: ", $(this).prop("class"));
    console.log("Button Text is: ", $(this).text(), $(this).text().length);
  });
});
