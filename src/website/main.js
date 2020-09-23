$(document).ready(function () {
  console.log("This is on start up");
  const url = "127.0.0.1:8080";
  var movesDisplayed = false;
  var moveDisplayedPiece = "00";
  var htmlChangedPiece = [];
  var htmlChangedPieceOrigImage = [];
  var whiteTurn = true
  getCurrentPlayerTurn();
  var promotedPawnLocation = ""
  var promotedPawnColour = ""
  const piecePixel = "100px 100px";

  const whitePieces = document.getElementsByName("white");
  const blackPieces = document.getElementsByName("black");
  for (let i = 0; i < whitePieces.length; i++) {
    setImage(
      whitePieces.item(i),
      "url('/static/imgs/" +
        whitePieces.item(i).value.toLowerCase() +
        "white" +
        whitePieces.item(i).getAttribute("class")[0] +
        ".png')"
    );
  }

  for (let i = 0; i < blackPieces.length; i++) {
    setImage(
      blackPieces.item(i),
      "url('static/imgs/" +
        blackPieces.item(i).value.toLowerCase() +
        "black" +
        blackPieces.item(i).getAttribute("class")[0] +
        ".png')"
    );
  }

  function swapTurn() {
    whiteTurn = !whiteTurn;
    setPlayerText()
  }

  function setPlayerText() {
    if (whiteTurn) {
        document.getElementById("playerText").innerHTML = "White's Turn";
      } else {
        document.getElementById("playerText").innerHTML = "Black's Turn";
      }
  }

  function checkText(set) {
    if (set) {
      document.getElementById("playerText").innerHTML =
        document.getElementById("playerText").innerHTML +
        "\tCHECK";
    } else {
      document.getElementById("playerText").innerHTML =
        document.getElementById("playerText").innerHTML + "";
    }
  }

  function mateText(set) {
    if (set) {
      document.getElementById("playerText").innerHTML =
        document.getElementById("playerText").innerHTML +
        "\tCHECKMATE\tGame Over";
    } else {
      document.getElementById("playerText").innerHTML =
        document.getElementById("playerText").innerHTML + "";
    }
    return set;
  }

  function setImage(object, url) {
    object.style.backgroundImage = url;
    object.style.backgroundSize = piecePixel;
    object.style.backgroundPosition = "center";
  }

  function setEmptyImage(id) {
    const pawnPiece = document.getElementById(id);
    pawnPiece.setAttribute("name", "empty");
    pawnPiece.setAttribute("value", " ");
    setImage(
        pawnPiece,
        "url('/static/imgs/" + pawnPiece.getAttribute("class") + ".png')"
      );
  }

  function getForm(message, name) {
    var myForm = document.createElement("FORM");
    myForm.setAttribute("method", "POST");
    var input = document.createElement("INPUT");
    input.setAttribute("type", "text");
    input.setAttribute("name", name);
    input.setAttribute("value", message);
    myForm.appendChild(input);
    return myForm
  }

  function getCurrentPlayerTurn() {
    fetch("/game", {
        method: "POST",
        body: new FormData(getForm("ply ", "empty")),
      })
        .then((response) => response.text())
        .then((data) => handleResponse(data, "ply", $(this)))
        .catch((error) => console.error("Error encountered: ", error));
  }

  function movePiece(newPieceID, oldPieceID) {
    const newPiece = document.getElementById(newPieceID);
    const oldPiece = document.getElementById(oldPieceID);
    newPiece.setAttribute("name", oldPiece.getAttribute("name"));
    newPiece.setAttribute("value", oldPiece.getAttribute("value"));
    setImage(
      newPiece,
      "url('/static/imgs/" +
        newPiece.value.toLowerCase() +
        newPiece.getAttribute("name") +
        newPiece.getAttribute("class")[0] +
        ".png')"
    );

    oldPiece.setAttribute("name", "empty");
    oldPiece.setAttribute("value", " ");
    setImage(
      oldPiece,
      "url('/static/imgs/" + oldPiece.getAttribute("class") + ".png')"
    );
  }

  function handleResponse(response, mode, clickedButton) {
    if (mode == "opt") {
      console.log("Operating mode: opt");
      console.log("Reponse from server:", response, "length:", response.length);
      for (let i = 0; i < response.length; i += 2) {
        const htmlPiece = document.getElementById(response.substring(i, i + 2));
        console.log(
          "Pushed to array ",
          response.substring(i, i + 2),
          htmlPiece.style.backgroundImage
        );
        htmlChangedPiece.push(response.substring(i, i + 2));
        htmlChangedPieceOrigImage.push(htmlPiece.style.backgroundImage);
        console.log(
          "THIS IS THE NEW IMAGE PATH:",
          htmlPiece.style.backgroundImage
        );
        if (htmlPiece.getAttribute("value") == " ") {
          setImage(
            htmlPiece,
            "url('/static/imgs/" + htmlPiece.getAttribute("class") + "s.png')"
          );
        } else {
          console.log(
            "THIS IS THE SUBSTRING:",
            htmlPiece.style.backgroundImage.substring(
              0,
              htmlPiece.style.backgroundImage.length - 6
            ) +
              "s" +
              '.png")'
          );
          setImage(
            htmlPiece,
            htmlPiece.style.backgroundImage.substring(
              0,
              htmlPiece.style.backgroundImage.length - 6
            ) +
              "s" +
              '.png")'
          );
        }
      }

      movesDisplayed = true;
      moveDisplayedPiece = clickedButton.prop("id");
    } else if (mode == "mov") {
      console.log("Operating mode: mov");
      console.log("Reponse from server:", response, "length:", response.length);
      var result = response.split(":")[1];
      if (result.substring(0, 4) === "true") {
        clearDisplayedMoves();
        movePiece(clickedButton.prop("id"), moveDisplayedPiece);
        swapTurn();
      } else if (result.substring(0, 6) === "castle") {
        console.log(
          "THIS IS THE ROOK:",
          result.substring(6, 8),
          result.substring(8, 10)
        );
        clearDisplayedMoves();
        movePiece(clickedButton.prop("id"), moveDisplayedPiece);
        movePiece(result.substring(8, 10), result.substring(6, 8));
        swapTurn();
      } else if (result.substring(0, 5) === "false") {
        console.log("THAT WAS AN INVALID MOVE");
      } else if (result.substring(0, 9) === "enpassant") {
        clearDisplayedMoves();
        movePiece(clickedButton.prop("id"), moveDisplayedPiece);
        setEmptyImage(result.substring(9, 11))
        swapTurn();
      } else if (result.substring(0, 3) == "pwn") {
          promotedPawnLocation = clickedButton.prop("id")
          promotedPawnColour = whiteTurn ? "white" : "black"
          document.getElementById("qPromoteTo").style.backgroundImage = "url('/static/imgs/q" + promotedPawnColour + "l.png')"
          document.getElementById("bPromoteTo").style.backgroundImage = "url('/static/imgs/b" + promotedPawnColour + "l.png')"
          document.getElementById("hPromoteTo").style.backgroundImage = "url('/static/imgs/h" + promotedPawnColour + "l.png')"
          document.getElementById("rPromoteTo").style.backgroundImage = "url('/static/imgs/r" + promotedPawnColour + "l.png')"
          $("#piecePromoteModal").modal('show')
          console.log("SHOWN MODAL")
      }
      console.log(
        "CHECK TEXT: ",
        result.substring(result.length - 5, result.length),
        result.substring(result.length - 5, result.length) === "check"
      );
      if (mateText(result.substring(result.length - 4, result.length) === "mate")) {
        return;
      }
      checkText(result.substring(result.length - 5, result.length) === "check");
    } else if (mode == "pwn") {
        if (response.substring(0, 4) == "true") {
            clearDisplayedMoves();
            console.log("Server: ", response, response.substring(response.length - 5, response.length))
            const newPiece = document.getElementById(promotedPawnLocation);
            newPiece.setAttribute("name", promotedPawnColour);
            newPiece.setAttribute("value", clickedButton.prop("id")[0]);
            setImage(
              newPiece,
              "url('/static/imgs/" +
                newPiece.value.toLowerCase() +
                newPiece.getAttribute("name") +
                newPiece.getAttribute("class")[0] +
                ".png')"
            );
            setEmptyImage(moveDisplayedPiece)
            if (mateText(response.substring(response.length - 4, response.length) === "mate")) {
                return;
            }
            swapTurn();
            checkText(response.substring(response.length - 5, response.length) === "check");
        } else {
            console.log("ERROR HAS OCCURED, PIECE CANNOT BE PROMOTED")
        }
    } else if (mode == "rst") {
        if (response.substring(0, 6) == "reload") {
            location.reload()
        }
    } else if (mode == "ply") {
        whiteTurn = (response.substring(0, 4) === "true" ? true : false)
        setPlayerText()
        if (mateText(response.substring(response.length - 4, response.length) === "mate")) {
            return;
        }
        checkText(response.substring(response.length - 5, response.length) === "check");
    }
  }

  function clearDisplayedMoves() {
    for (let i = 0; i < htmlChangedPiece.length; i++) {
      const htmlPiece = document.getElementById(htmlChangedPiece[i]);
      htmlPiece.style.backgroundImage = htmlChangedPieceOrigImage[i];
    }
    htmlChangedPiece = [];
    htmlChangedPieceOrigImage = [];
    movesDisplayed = false;
  }

  $("#return").click(function () {
    fetch("/game", {
        method: "POST",
        body: new FormData(getForm("rst ", "empty")),
      })
        .then((response) => response.text())
        .then((data) => handleResponse(data, "rst", $(this)))
        .catch((error) => console.error("Error encountered: ", error));
  });

  $(".piecePromobutton").click(function () {
    // var myForm = document.createElement("FORM");
    // myForm.setAttribute("method", "POST");
    // var input = document.createElement("INPUT");
    // input.setAttribute("type", "text");
    // input.setAttribute("name", "empty");
    // var mode = "pwn"
    // input.setAttribute("value", mode + " " + promotedPawnLocation + $(this).prop("id")[0]);
    // console.log(mode + " " + promotedPawnLocation + $(this).prop("id")[0])
    // console.log("PROMOTED THE PAWN AT", promotedPawnLocation, "TO A", $(this).prop("id")[0]);
    // myForm.appendChild(input);
    fetch("/game", {
        method: "POST",
        body: new FormData(getForm("pwn" + " " + promotedPawnLocation + $(this).prop("id")[0], "empty")),
      })
        .then((response) => response.text())
        .then((data) => handleResponse(data, "pwn", $(this)))
        .catch((error) => console.error("Error encountered: ", error));
  });

  $(".light, .dark").click(function () {
    // var myForm = document.createElement("FORM");
    // myForm.setAttribute("method", "POST");

    // var input = document.createElement("INPUT");
    // input.setAttribute("type", "text");
    // input.setAttribute("name", $(this).prop("name"));

    if (movesDisplayed) {
      if (htmlChangedPiece.includes($(this).prop("id"))) {
        console.log("Sending move request!!");
        //move the piece there
        var mode = "mov";
        // input.setAttribute(
        //   "value",
        //   mode + " " + moveDisplayedPiece + $(this).prop("id")
        // );
        var value = mode + " " + moveDisplayedPiece + $(this).prop("id") 
      } else {
        //Clear displayed moves
        console.log("These are the changed pieces", htmlChangedPiece);
        clearDisplayedMoves();
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
      var mode = "opt";
    //   input.setAttribute("value", mode + " " + $(this).prop("id"));
      var value = mode + " " + $(this).prop("id")
    }
    // myForm.appendChild(input);

    fetch("/game", {
      method: "POST",
      body: new FormData(getForm(value, $(this).prop("name"))),
    })
      .then((response) => response.text())
      .then((data) => handleResponse(data, mode, $(this)))
      .catch((error) => console.error("Error encountered: ", error));
  });

  $("button").click(function () {
    console.log("Value is: ", $(this).prop("id"));
    console.log("Name is: ", $(this).prop("name"));
    console.log("Class is: ", $(this).prop("class"));
    console.log("Button Text is: ", $(this).text(), $(this).text().length);
  });
});

