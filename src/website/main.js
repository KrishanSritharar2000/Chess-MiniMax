$(document).ready(function () {
  console.log("This is on start up");
  const url = "127.0.0.1:8080";
  var movesDisplayed = false;
  var moveDisplayedPiece = "00";
  var htmlChangedPiece = [];
  var htmlChangedPieceOrigImage = [];
  var whiteTurn = true
  getCurrentPlayerTurn();
  var thisIsWhitePlayer = true
  var calledMovFromOpp = false
  var calledFromBck = false
  var gameMode = 0
  var movesMade = 0
  var originalTakeBackMessage = document.getElementById("takeBackText").innerHTML
  var originalTakeBackMessageColour = document.getElementById("takeBackText").style.color
  getPlayerColourAndMode()
  console.log("THIS IS THE PLAYER COLOUR:", thisIsWhitePlayer)
  var promotedPawnLocation = ""
  var promotedPawnColour = ""
  var otherPieceClicked = 0
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
    document.getElementById("moves").innerHTML = "Moves:  " + movesMade
  }

  function checkText(set) {
    if (set) {
      document.getElementById("playerText").innerHTML =
        document.getElementById("playerText").innerHTML +
        "\tCHECK";
        document.getElementById("playerText")
        $("#playerText").addClass("animated pulse fast").one(
          "webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend",
          function () {
            $(this).removeClass("animated pulse fast");
          }
        );
    } else {
      document.getElementById("playerText").innerHTML =
        document.getElementById("playerText").innerHTML + "";
    }
  }

  // $("#playerText").one("animationend", document.getElementById("playerText").removeClass('animated pulse'))

  function mateText(set) {
    if (set) {
      document.getElementById("playerText").innerHTML =
        document.getElementById("playerText").innerHTML +
        "\tCHECKMATE\tGame Over";
        document.getElementById("playerText")
        $("#playerText").addClass("animated tada").one(
          "webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend",
          function () {
            $(this).removeClass("animated tada");
          }
        );
    } else {
      document.getElementById("playerText").innerHTML =
        document.getElementById("playerText").innerHTML + "";
    }
    return set;
  }

  function staleText(set) {
    if (set) {
      document.getElementById("playerText").innerHTML =
        document.getElementById("playerText").innerHTML +
        "\tSTALEMATE\tGame Over";
        $("#playerText").addClass("animated tada").one(
          "webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend",
          function () {
            $(this).removeClass("animated tada");
          }
        );
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

  function makeFetch(message, name, mode, clickedButton) {
    var myForm = document.createElement("FORM");
    myForm.setAttribute("method", "POST");
    var input = document.createElement("INPUT");
    input.setAttribute("type", "text");
    input.setAttribute("name", name);
    input.setAttribute("value", message);
    myForm.appendChild(input);
    fetch("/game", {
        method: "POST",
        body: new FormData(myForm),
      })
        .then((response) => response.text())
        .then((data) => handleResponse(data, mode, clickedButton))
        .catch((error) => console.error("Error encountered: ", error));
  }

  function getCurrentPlayerTurn() {
      makeFetch("ply ", "empty", "ply", $(this))
    // fetch("/game", {
    //     method: "POST",
    //     body: new FormData(getForm("ply ", "empty")),
    //   })
    //     .then((response) => response.text())
    //     .then((data) => handleResponse(data, "ply", $(this)))
    //     .catch((error) => console.error("Error encountered: ", error));
  }

  function getPlayerColourAndMode() {
    makeFetch("col", "empty", "col", $(this))
  }

  function getOpponentMove() {
    makeFetch("opp", "empty", "opp", $(this))
  }

  function getAIMove() {
    makeFetch("aim", "empty", "opp", $(this))
  }

  function isLetter(char){
    return ( (char >= 'A' &&  char <= 'Z') ||
             (char >= 'a' &&  char <= 'z') );
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
      "url('/static/imgs/" + oldPiece.getAttribute("class").substring(0, 5).trim() + ".png')"
    );
  }

  function handleResponse(response, mode, clickedButton) {
    if (mode == "opt") {
      console.log("Operating mode: opt");
      console.log("Reponse from server for opt:", response, "length:", response.length);
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
      console.log("Reponse from server mov:", response, "length:", response.length);
      if (document.getElementById("takeBackText").style.visibility == "visible") {
        document.getElementById("takeBackText").innerHTML = originalTakeBackMessage
        document.getElementById("takeBackText").style.color = originalTakeBackMessageColour
        document.getElementById("takeBackText").style.visibility = "hidden"
      }
      var result = response.split(":")[1];
      if (result.substring(0, 4) === "true") {
        clearDisplayedMoves();
        console.log("Before calling move piece:", clickedButton.prop("id"), moveDisplayedPiece)
        movePiece(clickedButton.prop("id"), moveDisplayedPiece);
        movesMade++
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
        return
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
          if (gameMode == 2) {
            if (thisIsWhitePlayer == whiteTurn) {
              $("#piecePromoteModal").modal('show')
              console.log("SHOWN MODAL")
            }
          } else {
            $("#piecePromoteModal").modal('show')
            console.log("SHOWN MODAL")
          }
      }
      console.log(
        "CHECK TEXT: ",
        result.substring(result.length - 5, result.length),
        result.substring(result.length - 5, result.length) === "check"
      );
      if (staleText(result.substring(result.length - 5, result.length) === "stale")) {
        return;
      }
      if (mateText(result.substring(result.length - 4, result.length) === "mate")) {
        return;
      }
      checkText(result.substring(result.length - 5, result.length) === "check");
      if (!calledMovFromOpp) {
        if (gameMode == 2) {
          getOpponentMove()
        } else if (gameMode == 1) {
          getAIMove()
        }
      }
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
            if (staleText(response.substring(response.length - 5, response.length) === "stale")) {
                return;
              }
            if (mateText(response.substring(response.length - 4, response.length) === "mate")) {
                return;
            }
            swapTurn();
            checkText(response.substring(response.length - 5, response.length) === "check");
        } else {
            console.log("ERROR HAS OCCURED, PIECE CANNOT BE PROMOTED")
        }
        if (gameMode == 2) {
          getOpponentMove()
        }
    } else if (mode == "rst") {
        if (response.substring(0, 6) == "reload") {
            location.reload()
        }
    } else if (mode == "ply") {
        whiteTurn = (response.substring(0, 4) === "true" ? true : false)
        index = whiteTurn ? 4 : 5
        console.log(index, response.charAt(index))
        while (!isLetter(response.charAt(index)) && index <= response.length) {
          console.log(response.charAt(index))
          index++
        }
        console.log(index)

        movesMade = parseInt(response.substring(whiteTurn ? 4 : 5, index-1 == response.length ? index + 1 : index))
        console.log("Moves ", response.substring(whiteTurn ? 4 : 5, index), " Made:", movesMade, response.substring(4,5), index)
        setPlayerText()
        if (staleText(response.includes("stale"))) {
            return;
          }
        if (mateText(response.includes("mate"))) {
            return;
        }
        checkText(response.includes("check"));
        if (thisIsWhitePlayer != whiteTurn && gameMode == 2) {
          getOpponentMove()
        }
    } else if (mode == "bck") {
        console.log("SERVER bck:", response)
        if (response.substring(0, 4) === "true") {
            movesMade--
            console.log("coord",response.substring(4,5), response.substring(5,6), response.substring(8,9) + response.substring(9,10))
            movePiece(response.substring(4,5) + response.substring(5,6), response.substring(8,9) + response.substring(9,10))
            const newPiece = document.getElementById(response.substring(8,9) + response.substring(9,10));
            var name = response.substring(11, 12) == "t" ? "black" : "white"
            name = response.substring(10, 11) == " " ? "empty" : name
            newPiece.setAttribute("name", name);
            newPiece.setAttribute("value", response.substring(10, 11));
            if (response.substring(10, 11) == " ") {
                setEmptyImage(response.substring(8,9) + response.substring(9,10))
            } else {
                setImage(
                newPiece,
                "url('/static/imgs/" + newPiece.value.toLowerCase() + newPiece.getAttribute("name") + newPiece.getAttribute("class")[0] + ".png')");
            }

            //check undoing a castle
            if (response.substring(6,7) == "K" && Math.abs(parseInt(response.substring(5,6)) - parseInt(response.substring(9,10))) == 2) {
              var rookLoc = ""
              var spaceLoc = ""
              if (response.substring(4,5) == "7") {
                //black king castled
                rookLoc = (response.substring(9,10) == "2") ? "70" : "77"
                spaceLoc = (response.substring(9,10) == "2") ? "73" : "75"
              } else {
                //white king castled
                rookLoc = (response.substring(9,10) == "2") ? "00" : "07"
                spaceLoc = (response.substring(9,10) == "2") ? "03" : "05"
              }
              const rookPiece = document.getElementById(rookLoc);
              rookPiece.setAttribute("name",  response.substring(7, 8) == "t" ? "black" : "white");
              rookPiece.setAttribute("value", "R");
              setImage(rookPiece, "url('/static/imgs/" + rookPiece.value.toLowerCase() + rookPiece.getAttribute("name") + rookPiece.getAttribute("class")[0] + ".png')")
              setEmptyImage(spaceLoc)
            }

            //check undoing pawn promotion
            if (response.substring(6,7) == "P" && (response.substring(8,9) == "7" || response.substring(8,9) == "0")) {
              var pawnPiece = document.getElementById(response.substring(4,5) + response.substring(5,6));
              pawnPiece.setAttribute("value", "p")
              setImage(
                pawnPiece, "url('/static/imgs/" + pawnPiece.value.toLowerCase() + pawnPiece.getAttribute("name") + pawnPiece.getAttribute("class")[0] + ".png')");
            }

            swapTurn()
        if (response.substring(response.length - 2, response.length) == "ai") {
          handleResponse("true"+response.substring(response.length - 10, response.length - 2), "bck", clickedButton)
        } else if (response.substring(response.length - 2, response.length) == "ac") {
          if (document.getElementById("takeBackText").style.visibility == "visible") {
            document.getElementById("takeBackText").innerHTML = "Takeback accepted by opponent<i class='fas fa-check ml-2 text-success'></i>"
            document.getElementById("takeBackText").style.color = "green"
          }
          calledFromBck = true
          handleResponse("true"+response.substring(response.length - 10, response.length - 2), "bck", clickedButton)
          calledFromBck = false
          console.log("TAKEBACK ACCEPTED")
        }
      } else if (response.substring(0, 6) == "reject") {
        console.log("TAKEBACK REJECTED")
        if (document.getElementById("takeBackText").style.visibility == "visible") {
          document.getElementById("takeBackText").innerHTML = "Takeback rejected by opponent<i class='fas fa-times ml-2 text-danger'></i>"
          document.getElementById("takeBackText").style.color = "#df4759"
        }
      } else {
          document.getElementById("takeBackText").innerHTML = "No moves left to Takeback<i class='fas fa-times ml-2 text-danger'></i>"
          document.getElementById("takeBackText").style.color = "white"
          document.getElementById("takeBackText").style.visibility = "visible"
      }
        if (gameMode == 2 && !calledFromBck && thisIsWhitePlayer != whiteTurn) {
          getOpponentMove()
        }
      
        // if (!thisIsWhitePlayer) {
        //   getOpponentMove()
        // }
    } else if (mode == "col") {
      thisIsWhitePlayer = response.substring(0, 4) == "true"
      console.log(response, response.substring(0, 4), thisIsWhitePlayer)
      gameMode = parseInt(response.slice(response.length - 1))
      console.log("THIS IS THE GAME MODE", gameMode, response.slice(response.length - 1))
      if (gameMode == 2 && thisIsWhitePlayer != whiteTurn) {
        getOpponentMove()
      } else if (gameMode == 1 && !thisIsWhitePlayer) {
        getAIMove()
      }
    } else if (mode == "opp") {
      if (response.substring(0,3) == "bck") {
        console.log("TAKE BACK MODAL SHOWN")
        $("#takeBackMoveModal").modal('show')
      } else {
        console.log("Server opp:", response)
        var respLen = response.length
        moveDisplayedPiece = response.substring(respLen - 4, respLen - 2)
        if (document.getElementById(moveDisplayedPiece).getAttribute("value") == " ") {
          console.log("Invalid opp request")
          return
        }
        var movedTo = response.substring(respLen - 2, respLen)
        console.log("MoveDisplayedPiece:", response.substring(respLen - 4, respLen - 2), "MovedTo", response.substring(respLen - 2, respLen))
        console.log("Passing on response:", response.substring(0, respLen - 4))
        calledMovFromOpp = true
        handleResponse(response.substring(0, respLen - 4), "mov", $('#' + movedTo))
        calledMovFromOpp = false
      }
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

  $("#restart").click(function () {
    makeFetch("rst ", "empty", "rst", $(this))
    // fetch("/game", {
    //     method: "POST",
    //     body: new FormData(getForm("rst ", "empty")),
    //   })
    //     .then((response) => response.text())
    //     .then((data) => handleResponse(data, "rst", $(this)))
    //     .catch((error) => console.error("Error encountered: ", error));
  });

  $("#undo").click(function () {
    if (gameMode == 2) {
      if (thisIsWhitePlayer == whiteTurn) {
        if (movesMade < 2) {
          document.getElementById("takeBackText").innerHTML = "Cannot undo any more moves<i class='fas fa-times ml-2 text-white'></i>"
          document.getElementById("takeBackText").style.color = originalTakeBackMessageColour
          document.getElementById("takeBackText").style.visibility = "visible"
          return
        }
        if (document.getElementById("takeBackText").innerHTML != originalTakeBackMessage) {
          document.getElementById("takeBackText").innerHTML = originalTakeBackMessage
          document.getElementById("takeBackText").style.color = originalTakeBackMessageColour
        }
        document.getElementById("takeBackText").style.visibility = "visible"
      } else {
        document.getElementById("takeBackText").innerHTML = "Can only undo move on your turn<i class='fas fa-times-circle ml-2 text-white'></i>"
        document.getElementById("takeBackText").style.color = originalTakeBackMessageColour
        document.getElementById("takeBackText").style.visibility = "visible"
        return
      }
    }
    makeFetch("bck ", "empty", "bck", $(this))
  });

  $("#acceptTakeBack").click(function () {
    makeFetch("atb ", "empty", "bck", $(this))
  })

  $("#rejectTakeBack").click(function () {
    makeFetch("rtb ", "empty", "bck", $(this))
  })

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
    // fetch("/game", {
    //     method: "POST",
    //     body: new FormData(getForm("pwn" + " " + promotedPawnLocation + $(this).prop("id")[0], "empty")),
    //   })
    //     .then((response) => response.text())
    //     .then((data) => handleResponse(data, "pwn", $(this)))
    //     .catch((error) => console.error("Error encountered: ", error));
    makeFetch("pwn" + " " + promotedPawnLocation + $(this).prop("id")[0], "empty", "pwn", $(this))
  });

  $(".light, .dark").click(function () {
    $(this).addClass("animated pulse faster").one(
      "webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend",
      function () {
        $(this).removeClass("animated pulse faster");
      }
    );
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
      var turnChecker = whiteTurn
      if (gameMode == 2) {
        turnChecker = thisIsWhitePlayer
      }
      if (
        (turnChecker && $(this).prop("name") == "black") ||
        (!turnChecker && $(this).prop("name") == "white")
      ) {
        console.log("NO REQUEST TO SERVER. THAT IS OPPONENTS PIECE");
        if (++otherPieceClicked >= 2) {
          $("#playerText").addClass("animated flash").one(
            "webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend",
            function () {
              $(this).removeClass("animated flash");})
          otherPieceClicked = 0
        }
        return;
      }
      if ($(this).prop("name") == "empty") {
        console.log("NO REQUEST TO SERVER. THERE IS NO PIECE THERE");
        return;
      }
      otherPieceClicked = 0
      // input.setAttribute("value", $(this).prop("name") + " " + $(this).val() + " piece selected");
      var mode = "opt";
    //   input.setAttribute("value", mode + " " + $(this).prop("id"));
      var value = mode + " " + $(this).prop("id")
    }
    // myForm.appendChild(input);

    // fetch("/game", {
    //   method: "POST",
    //   body: new FormData(getForm(value, $(this).prop("name"))),
    // })
    //   .then((response) => response.text())
    //   .then((data) => handleResponse(data, mode, $(this)))
    //   .catch((error) => console.error("Error encountered: ", error));
    makeFetch(value, $(this).prop("name"), mode, $(this))
  });

  $("button").click(function () {
    console.log("Value is: ", $(this).prop("id"));
    console.log("Name is: ", $(this).prop("name"));
    console.log("Class is: ", $(this).prop("class"));
    console.log("Button Text is: ", $(this).text(), $(this).text().length); 
  });
});

