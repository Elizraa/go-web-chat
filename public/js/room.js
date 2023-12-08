function loginUserRequest(username, password, resolve = console.log, reject = console.log) {
    console.log("logging in", username);
    $.ajax({
        url: '/auth',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({ username: username, password: password }),
        success: function (data, textStatus, xhr) {
            try {
                // Check if the response is a valid JSON
                if (typeof data === 'object' && data !== null && data.status) {
                    resolve(data);
                } else {
                    // Log the entire response for debugging
                    console.log("Invalid response:", data);

                    // Reject with an appropriate message
                    reject("Invalid response for login request");
                }
            } catch (error) {
                console.error("Error processing the response:", error);

                // Reject with an appropriate message
                reject("Error processing the response for login request");
            }
        },
        error: function (xhr, status, error) {
            console.log("Status:", status);
            console.log("Error:", error);
            reject("Could not log in " + username);
        }
    });
}


// REST API calls
// POST /chats
function createRoom(title, description, visibility, password, resolve=console.log, reject=console.log) {
    console.log("creating room " + title);
   $.post('/chats', JSON.stringify({ title: title, description: description, visibility: visibility, password: password }), "json")
       .done(function (data) {
           if (data.status) {
               console.log("successfully created room!")
               resolve(data)
           }
           else {
                console.log(data)
               reject("Could not create room " + title)
           }
       })
       .fail(function (xhr) {
            console.log(xhr)
           reject("Could not create room " + title)
       });
}

// GET /chats/<id>
function retrieveRoom(title,resolve=console.log,reject=console.log) {
   $.get('/chats/' + title)
       .done(function (data) {
           if (!data.hasOwnProperty('error')) {
               resolve(data);
           }
           else {
               reject(data);
           }
       })
       .fail(function (xhr) {
           reject("Error fetching chat room  " + title)
       });
}

// GET /chats/<id>
function retrieveRoomID(ID) {
   $.get('/chats/' + ID)
       .done(function (data) {
           console.log(data)
           if (!data.hasOwnProperty('error')) {
               console.log("success")
               return data;
           }
           else {
               console.warn("Could not retrieve chat room  " + ID);
           }
       })
       .fail(function (xhr) {
           console.log("Error fetching chat room " + ID);
           console.log(xhr);
       });
}
// PUT /chats/<id>
function putRoom(title, description, visibility, password) {
   $.ajax({
       url: "/chats/" + title,
       method: 'PUT',
       data: JSON.stringify({ title: title, description: description, visibility: visibility, password: password })
   })
       .done(function (data) {
           console.log(data)
           if (!data.hasOwnProperty('error')) {
               console.log("success")
               return data;
           }
           else {
               displayAlert("Could not update chat room  " + title);
           }
       })
       .fail(function (xhr) {
           console.log("Error fetching chat room " + title);
           console.log(xhr);
       });
}

// DELETE /chats/title
function deleteRoom(title) {
   $.ajax({
       url: "/chats/" + title,
       method: 'DELETE'
   })
       .done(function (data) {
           console.log(data)
           if (!data.hasOwnProperty('error')) {
               console.log("successfully deleted " + title)
               return data;
           }
           else {
               displayAlert("Could not delete chat room  " + title);
           }
       })
       .fail(function (xhr) {
           console.log("Error deleting chat room " + title);
           console.log(xhr);
       });
}

// DELETE /chats/ID
function deleteRoomID(ID) {
   $.ajax({
       url: "/chats/" + ID,
       method: 'DELETE'
   })
       .done(function (data) {
           console.log(data)
           if (!data.hasOwnProperty('error')) {
               console.log("successfully deleted " + ID)
               return data;
           }
           else {
               displayAlert("Could not delete chat room  " + ID);
           }
       })
       .fail(function (xhr) {
           console.log("Error deleting chat room " + ID);
           console.log(xhr);
       });
}