
$(function() {
    $("form[name='registration']").validate({

      rules: {
        name: {
          required: true,
          name: true
        },
        email: {
          required: true,
          email2: true,
        },
      
        postalcode: {
          required: true,
          postalcodeNL: true
        },
        
        address: {
          required: true,
          address: true
        },
        phoneNumber:{
          required: true,
          phoneNL: true
        }
      },

      messages: {
        name: "Vul uw naam in!",
        email: "Vul een geldig mailadres in!",
        postalcode: "Vul een geldige postcode in!",
        address: "vul een adres in!"
      }, 
      submitHandler: function(form) {
        form.submit();
        alert('valid form');
            return false;
      },
    });
  }); 