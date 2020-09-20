$(document).ready(function(){
    console.log("This is on start up")
    const url = "127.0.0.1:8080"
    var movesDisplayed = false
    var moveDisplayedPiece = "00"
    var htmlPiecesWithDisplayedMove = []
    var whiteTurn = true

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

    function handleResponse(response, mode, clickedButton) {
      if (mode == "opt") {
        console.log("Operating mode: opt")
        console.log("Reponse from server:", response, "length:", response.length)
        for (let i = 0; i < response.length; i += 2) {
          const htmlPiece = document.getElementById(response.substring(i,i+2))
          console.log("Pushed to array ", response.substring(i,i+2), htmlPiece.innerHTML)
          htmlPiecesWithDisplayedMove.push(response.substring(i,i+2))
          htmlPiecesWithDisplayedMove.push(htmlPiece.innerHTML)
          htmlPiece.innerHTML = "TAKE " + htmlPiece.innerHTML
        }
        console.log("Array", htmlPiecesWithDisplayedMove, htmlPiecesWithDisplayedMove.length)
        console.log("Array", htmlPiecesWithDisplayedMove[0])
        console.log("Array", htmlPiecesWithDisplayedMove[1])


        movesDisplayed = true
        moveDisplayedPiece = clickedButton.prop("id")
      } else if (mode == "mov") {        
        console.log("Operating mode: mov")
        console.log("Reponse from server:", response, "length:", response.length)
      }

    }

    $(".light, .dark").click(function(){
      console.log("IS MOVES DISPLAYED:", movesDisplayed)
      var myForm = document.createElement("FORM");
      myForm.setAttribute("method", "POST");
      
      var input = document.createElement("INPUT");
      input.setAttribute("type", "text");
      input.setAttribute("name", $(this).prop("name"))

      var mode = "opt"
      if (movesDisplayed) {
        //Clear displayed moves
        console.log("These are the changed pieces", htmlPiecesWithDisplayedMove)
        for (let i = 0; i < htmlPiecesWithDisplayedMove.length; i += 2) {
          const htmlPiece = document.getElementById(htmlPiecesWithDisplayedMove[i])
          htmlPiece.innerHTML = htmlPiecesWithDisplayedMove[i+1]
        }
        htmlPiecesWithDisplayedMove = []
        movesDisplayed = false
        if ($(this).prop("id") == moveDisplayedPiece) {
          return
        }
      }

      if (!movesDisplayed) {
        if ((whiteTurn && $(this).prop("name") == "black") || (!whiteTurn && $(this).prop("name") == "white")) {
          console.log("NO REQUEST TO SERVER. THAT IS OPPONENTS PIECE")
          return
        }
        if ($(this).prop("name") == "empty") {
          console.log("NO REQUEST TO SERVER. THERE IS NO PIECE THERE")
          return
        }
        // input.setAttribute("value", $(this).prop("name") + " " + $(this).val() + " piece selected");
        input.setAttribute("value", mode + " " + $(this).val());
      }
      myForm.appendChild(input)

      fetch("/game", {
          method:"POST",
          body: new FormData(myForm)
      }).then(
        response => response.text()
      ).then(
        (data) => handleResponse(data, mode, $(this)) 
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