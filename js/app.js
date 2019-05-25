$(document).foundation()
function setConfirmMessage(confirm_password) {
 var password = document.getElementById("password").value;
 var message = "";
 if (password == confirm_password) {
   message = "";
 } else {
   message =  "パスワードが一致しません";
 }

 var div = document.getElementById("pass_confirm_message");
 if (!div.hasFistChild) {div.appendChild(document.createTextNode(""));}
 div.firstChild.data = message;
}
