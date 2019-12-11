$(document).ready(function(){
    $("input").on('input', function(){
      var msg = $(this).val();
      $("#Content").replaceWith("<div id=\"Content\">"+msg+"</div>");
    });
    $("input").change(function(){
        var msg = $(this).val();
        var key = $(this).attr('name');
      $("#Content").replaceWith("<div id=\"Content\">"+msg+"</div>");
      $.ajax({
        type: "POST",
        url: '/post',
        data: JSON.stringify({ redis_key: key, redis_val: msg }),
        contentType : 'application/json',
        dataType: 'json',
        success: function(data) {
            // alert(data.redis_key)
            
            $( "<p>"+data.redis_val+"</p>" ).appendTo( "#"+data.redis_key );
        }
      });
    });
  });