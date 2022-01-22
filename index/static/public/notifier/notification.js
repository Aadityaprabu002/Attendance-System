document.addEventListener('DOMContentLoaded', function() {
    if (!Notification) {
     alert('Desktop notifications not available in your browser. Try Different Browser.');
     return;
    }
   
    if (Notification.permission !== 'granted')
     Notification.requestPermission();
   });
   
   
   function notifyMe() {
    if (Notification.permission !== 'granted')
     Notification.requestPermission();
    else {
     var notification = new Notification('ATTENDANCE!', {
      icon: 'https://cdn.imgbin.com/9/23/23/imgbin-computer-icons-time-attendance-clocks-hourglass-icon-design-sd-card-83XQRXGwLkUwd4sdmnBrbiddv.jpg',
      body: 'Answer your attendance please!',
     });
     notification.onclick = function() {
      window.open('/student/dashboard/session');
     };
    }
   }