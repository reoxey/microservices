
function loadAddress(){
    AJAX("shipping/addresses", "GET", "", function(res, status, xhr){
        const obj = $('#all-addresses');
        if (xhr.status === 200) {
            obj.empty();
            for(let ob of res) {
                obj.append('<div class="col-10"><label class="checkbox-inline" for="addr-'+ob.id+'"><input name="address" id="addr-'+ob.id+'" type="radio" value="'+ob.id+'"> <span id="contact-name-'+ob.id+'">'+ob.contact_name+'</span>, <span id="landmark-'+ob.id+'">'+ob.city+'</span>, <span id="city-'+ob.id+'">'+ob.city+'</span>, <br><span class="state-'+ob.id+'">'+ob.state+'</span>, <span class="country-'+ob.id+'">'+ob.country+'</span>, <span class="zip-'+ob.id+'">'+ob.zip+'</span>, <br><span class="contact-phone-'+ob.id+'">'+ob.contact_phone+'</span></label></div><div class="col-2"><a href="#" onclick="editAddress('+ob.id+')"><i class="fa fa-pencil"></i></a>&nbsp;&nbsp;&nbsp;<a href="#" onclick="deleteAddress('+ob.id+')"><i class="fa fa-trash"></i></a></div>');
            }
        }

        $('#all-addresses').find('input[name="address"]').on('change', function(){
            $('input[name="address"]').parent("label").removeClass("checked");
            if($(this).is(':checked')){
                $(this).parent("label").addClass("checked");
            } else {
                $(this).parent("label").removeClass("checked");
            }
        });
    })
}

$('#address-part-2').slideUp();

$('#add-address').on('click', function(){
    $('#address-part-1').slideUp();
    $('#address-part-2').slideDown();
})

$('#address-form').on('submit', function(e){

    e.preventDefault();

    const id = $('#address-id').val();
    let method;
    if(id === "") method = "POST";
    else method = "PUT";

    const data = '{"contact_name": "'+$('#contact-name').val()+'",'+
    '"contact_phone": "'+$('#contact-phone').val()+'",'+
    '"landmark": "'+$('#landmark').val()+'",'+
    '"country": "'+$('#country').val()+'",'+
    '"state": "'+$('#state').val()+'",'+
    '"city": "'+$('#city').val()+'",'+
    '"zip": '+$('#zip-code').val()+'}';

    AJAX('shipping/addresses', method, data, function(res, status, xhr){
        
        if(xhr.status === 201 || xhr.status === 200) {
            loadAddress();
            $('#address-part-1').slideDown();
            $('#address-part-2').slideUp();
        } else {
            alertPop("danger", "Something went wrong. status: "+xhr.status)
        }
    })
    return false
});

function editAddress(id){
    //TODO
}

function deleteAddress(id){
    //TODO
}

$('#checkout').on('click', function(){

    const addr = $('input[name="address"]:checked').val();

    console.log(addr)

    if (!addr || addr === "") {
        alertPop("warning", "Please select delivery address.")
        return;
    }

    const data = '{"address_id": '+addr+','+
    '"payment_method": '+$('input[name="payment"]:checked').val()+'}';

    AJAX("cart/checkout", "POST", data, function(res, status, xhr){
        if(xhr.status === 200) {

            window.location = "./home.html";
        } else {
            alertPop("danger", "Something went wrong. status: "+xhr.status)
        }
    })
});

loadAddress();