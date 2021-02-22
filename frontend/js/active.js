/* =====================================
Template Name: Eshop
Author Name: Naimur Rahman
Author URI: http://www.wpthemesgrid.com/
Description: Eshop - eCommerce HTML5 Template.
Version:1.0
========================================*/

function setCookie(name,value,days) {
	var expires = "";
	if (days) {
		var date = new Date();
		date.setTime(date.getTime() + (days*24*60*60*1000));
		expires = "; expires=" + date.toUTCString();
	}
	document.cookie = name + "=" + (value || "")  + expires + "; path=/";
  }

function getCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
    for(var i=0;i < ca.length;i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1,c.length);
        if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
    }
    return null;
  }

  function eraseCookie(name) {   
    document.cookie = name +'=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
  }

let cartId = getCookie("cart_id");

(function($) {
    "use strict";
     $(document).on('ready', function() {	
		
		/*====================================
		03. Sticky Header JS
		======================================*/ 
		jQuery(window).on('scroll', function() {
			if ($(this).scrollTop() > 200) {
				$('.header').addClass("sticky");
			} else {
				$('.header').removeClass("sticky");
			}
		});
		
		/*=======================
		  Search JS JS
		=========================*/ 
		$('.top-search a').on( "click", function(){
			$('.search-top').toggleClass('active');
		});
		
		/*=======================
		  Slider Range JS
		=========================*/ 
		$( function() {
			$( "#slider-range" ).slider({
			  range: true,
			  min: 0,
			  max: 100,
			  values: [ 5, 50 ],
			  slide: function( event, ui ) {
				$( "#amount" ).val( "$" + ui.values[ 0 ] + " - $" + ui.values[ 1 ] );
			  }
			});
			$( "#amount" ).val( "$" + $( "#slider-range" ).slider( "values", 0 ) +
			  " - $" + $( "#slider-range" ).slider( "values", 1 ) );
		} );
		
		/*===========================
		  Quick View Slider JS
		=============================*/ 
		$('.quickview-slider-active').owlCarousel({
			items:1,
			autoplay:true,
			autoplayTimeout:5000,
			smartSpeed: 400,
			autoplayHoverPause:true,
			nav:true,
			loop:true,
			merge:true,
			dots:false,
			navText: ['<i class=" ti-arrow-left"></i>', '<i class=" ti-arrow-right"></i>'],
		});
		
		
		/*====================================
		  Cart Plus Minus Button
		======================================*/
		var CartPlusMinus = $('.cart-plus-minus');
		CartPlusMinus.prepend('<div class="dec qtybutton">-</div>');
		CartPlusMinus.append('<div class="inc qtybutton">+</div>');
		$(".qtybutton").on("click", function() {
			var $button = $(this);
			var oldValue = $button.parent().find("input").val();
			if ($button.text() === "+") {
				var newVal = parseFloat(oldValue) + 1;
			} else {
				// Don't allow decrementing below zero
				if (oldValue > 0) {
					var newVal = parseFloat(oldValue) - 1;
				} else {
					newVal = 1;
				}
			}
			$button.parent().find("input").val(newVal);
		});
		
		/*=======================
		  Extra Scroll JS
		=========================*/
		$('.scroll').on("click", function (e) {
			var anchor = $(this);
				$('html, body').stop().animate({
					scrollTop: $(anchor.attr('href')).offset().top - 0
				}, 900);
			e.preventDefault();
		});
		
		/*===============================
		10. Checkbox JS
		=================================*/  
		
		$('input[name="payment"]').change(function(){
			$('input[name="payment"]').parent("label").removeClass("checked");
			if($(this).is(':checked')){
				$(this).parent("label").addClass("checked");
			} else {
				$(this).parent("label").removeClass("checked");
			}
		});
		
		/*==================================
		 12. Product page Quantity Counter
		 ===================================*/
		$('.qty-box .quantity-right-plus').on('click', function () {
			var $qty = $('.qty-box .input-number');
			var currentVal = parseInt($qty.val(), 10);
			if (!isNaN(currentVal)) {
				$qty.val(currentVal + 1);
			}
		});
		$('.qty-box .quantity-left-minus').on('click', function () {
			var $qty = $('.qty-box .input-number');
			var currentVal = parseInt($qty.val(), 10);
			if (!isNaN(currentVal) && currentVal > 1) {
				$qty.val(currentVal - 1);
			}
		});
		
		
		/*====================================
			Scroll Up JS
		======================================*/
		$.scrollUp({
			scrollText: '<span><i class="fa fa-angle-up"></i></span>',
			easingType: 'easeInOutExpo',
			scrollSpeed: 900,
			animation: 'fade'
		});  
		
	});
	
	/*====================================
	18. Nice Select JS
	======================================*/	
	$('select').niceSelect();
		
	/*=====================================
	 Others JS
	======================================*/ 	
	$( function() {
		$( "#slider-range" ).slider({
			range: true,
			min: 0,
			max: 100,
			values: [ 10, 50 ],
			slide: function( event, ui ) {
				$( "#amount" ).val( "$" + ui.values[ 0 ] + " - $" + ui.values[ 1 ] );
			}
		});
		$( "#amount" ).val( "$" + $( "#slider-range" ).slider( "values", 0 ) +
		  " - $" + $( "#slider-range" ).slider( "values", 1 ) );
	} );
	
	/*=====================================
	  Preloader JS
	======================================*/ 	
	//After 2s preloader is fadeOut
	$('.preloader').delay(2000).fadeOut('slow');
	setTimeout(function() {
	//After 2s, the no-scroll class of the body will be removed
	$('body').removeClass('no-scroll');
	}, 2000); //Here you can change preloader time
	 
})(jQuery);

function bindBtnCart(inCart){
	//------------- DETAIL ADD - MINUS COUNT ORDER -------------//
$('.btn-number').click(function(e){
    e.preventDefault();
    
   let fieldName = $(this).attr('data-field');
    let type      = $(this).attr('data-type');
    var input = $("input[name='"+fieldName+"']");
    var currentVal = parseInt(input.val());
    if (!isNaN(currentVal)) {
        if(type == 'minus') {
            
            if(currentVal > input.attr('data-min')) {
                input.val(currentVal - 1).change();
            } 
            if(parseInt(input.val()) == input.attr('data-min')) {
                $(this).attr('disabled', true);
            }

        } else if(type == 'plus') {

            if(currentVal < input.attr('data-max')) {
                input.val(currentVal + 1).change();
            }
            if(parseInt(input.val()) == input.attr('data-max')) {
                $(this).attr('disabled', true);
            }

        }
    } else {
        input.val(0);
    }
	const pid = input.attr('data-pid');
	$('#modal-cart').attr('onclick', 'addToCart('+pid+', '+input.val()+')');

	if (inCart) triggerCartQtyTimeout(pid, input.val())
});
$('.input-number').focusin(function(){
   $(this).data('oldValue', $(this).val());
});
$('.input-number').change(function() {
    
    let minValue =  parseInt($(this).attr('data-min'));
    let maxValue =  parseInt($(this).attr('data-max'));
    let valueCurrent = parseInt($(this).val());
    
    let name = $(this).attr('name');
    if(valueCurrent >= minValue) {
        $(".btn-number[data-type='minus'][data-field='"+name+"']").removeAttr('disabled')
    } else {
        alert('Sorry, the minimum value was reached');
        $(this).val($(this).data('oldValue'));
    }
    if(valueCurrent <= maxValue) {
        $(".btn-number[data-type='plus'][data-field='"+name+"']").removeAttr('disabled')
    } else {
        alert('Sorry, the maximum value was reached');
        $(this).val($(this).data('oldValue'));
    }
	const pid = $(this).attr('data-pid');
    $('#modal-cart').attr('onclick', 'addToCart('+pid+', '+$(this).val()+')');
    if(inCart) triggerCartQtyTimeout(pid, $(this).val())
});
$(".input-number").keydown(function (e) {
        // Allow: backspace, delete, tab, escape, enter and .
        if ($.inArray(e.keyCode, [46, 8, 9, 27, 13, 190]) !== -1 ||
             // Allow: Ctrl+A
            (e.keyCode == 65 && e.ctrlKey === true) || 
             // Allow: home, end, left, right
            (e.keyCode >= 35 && e.keyCode <= 39)) {
                 // let it happen, don't do anything
                 return;
        }
        // Ensure that it is a number and stop the keypress
        if ((e.shiftKey || (e.keyCode < 48 || e.keyCode > 57)) && (e.keyCode < 96 || e.keyCode > 105)) {
            e.preventDefault();
        }
    });
}
bindBtnCart(false);
let timer = null;
function triggerCartQtyTimeout(pid, qty){
	  
		if(timer) {
		  clearTimeout(timer)
		}
	  
		timer = setTimeout(()=>{
			addToCart(pid, qty, "PUT")
		}, 2000);
}

  const auth = getCookie("jwt_auth");

  console.log(auth);
  if (auth === null) {
      window.location = "./index.html"
  }

  $('#log-out').on('click', function(){
      eraseCookie("jwt_auth");
	  eraseCookie("cart_id")
      window.location = "./index.html"
  });

  function AJAX(path, method, data, fn) {
    $.ajax({
        url: 'http://micro.reoxey.com/'+path,
        method: method,
        headers: {Authorization: "Bearer "+auth },
        contentType: "application/json",
        data: data,
        dataType: "json",
        success: fn,
        error: function(err){
          alertPop("danger", "Something went wrong!");
		  console.error(err)
        }
    });
  }

  function modalData(id, stocks) {
      $('#modal-title').text($('#product-title-'+id).text());
	  const modalStocks = $('#modal-stocks');
      if (stocks <= 0) {
        modalStocks.html('<span><b class="text-danger">Out of stock</b></span>');
      } else 
      if (stocks < 5) {
        modalStocks.html('<span><b class="text-danger">Only '+stocks+' remaining</b></span>');  
      } else {
        modalStocks.html('<span><i class="fa fa-check-circle-o"></i> in stock</span>');
      }
	  const productStocks = $('#product-stocks');
	  productStocks.val(1);
	  productStocks.attr('data-max', stocks);
	  productStocks.attr('data-pid', id);
      $('#modal-price').text($('#product-price-'+id).text());
      $('#modal-cart').attr('onclick', 'addToCart('+id+', 1)');
      $('#exampleModal').modal('show');
  }

  function alertPop(t, d){
	  $('#alert-here').html('<div class="alert alert-'+t+' alert-dismissible fade show" role="alert">'+
	  '<strong>Alert!</strong> '+ d+
	  '<button type="button" class="close" data-dismiss="alert" aria-label="Close">'+
		'<span aria-hidden="true">&times;</span>'+
	  '</button></div>')
  }

  console.log(cartId);
  if (cartId === null) {
    
	AJAX("cart", "POST", "", function(res, status, xhr){
		if (xhr.status === 201){
			const loc = xhr.getResponseHeader('Location');
			if (loc === null) return

			cartId = loc.split("/")[3]

			setCookie("cart_id", cartId, 30)
		}
	})

  } else {
	  loadCartItems()
  }

  function round(x){
	  return Number(x).toFixed(3)
  }

  function loadCartItems(){
	AJAX("cart/"+cartId, "GET", "", function(res, status, xhr){
		if (xhr.status === 200) {
			const mainCart = $('#shopping-cart-area');
			const obj = $('#shopping-cart')
			console.log(res)
			if(res.items === null) {
				if(mainCart.length) mainCart.html('<section class="shop-services section home"><div class="container"><div class="row"><div class="col-12"><div class="single-service"><h5 class="mb-2">Your Cart is empty</h5> <a href="./home.html" class="btn pull-right text-light">Continue shopping</a></div></div></div></div> </section>');
				obj.html('');
				$('#total-count').html('<i class="fa fa-shopping-cart"></i>');
				return;
			}
			let cartItems = "";
			let cartCount = 0;
			for (let ob of res.items) {
				cartCount++;
				cartItems += '<li>'+
				'<a href="#" class="remove" title="Remove this item" onclick="removeFromCart('+ob.id+')"><i class="fa fa-remove"></i></a>'+
				'<a class="cart-img" href="#"><img src="https://via.placeholder.com/70x70" alt="#"></a>'+
				'<h4><a href="#">'+ob.name+'</a></h4>'+
				'<p class="quantity">'+ob.qty+'x - <span class="amount">$'+round(ob.qty*ob.price)+'</span></p>'+
			'</li>';
			}
			$('#total-count').html('<i class="fa fa-shopping-cart"></i><span class="total-count">'+cartCount+'</span>');
			obj.html('<div class="dropdown-cart-header">'+
			'<span>'+cartCount+' Items</span>'+
			'<a href="./cart.html">View Cart</a></div>'+
			'<ul class="shopping-list">' + cartItems + '</ul>'+
			'<div class="bottom">'+
			'<div class="total"><span>Total</span>'+
			'<span class="total-amount">$'+round(res.payment.total)+'</span></div>'+
			'<a href="checkout.html" class="btn animate">Checkout</a></div>'
			)

			if (mainCart.length) {
				let html = '<table class="table shopping-summery"><thead><tr class="main-hading"><th>PRODUCT</th><th>NAME</th><th class="text-center">UNIT PRICE</th><th class="text-center">QUANTITY</th><th class="text-center">TOTAL</th><th class="text-center"><i class="ti-trash remove-icon"></i></th></tr></thead><tbody>';

				for(let ob of res.items) {
					let old_price = '';
					if (ob.old_price){
						old_price = '<span class="oldprice">$'+ob.old_price+'</span><br>';
					}

					html += '<tr><td class="image" data-title="No"><img src="https://via.placeholder.com/100x100" alt="#"></td><td class="product-des" data-title="Description"><p class="product-name"><a href="#">'+ob.name+'</a></p><p class="product-des">Maboriosam in a tonto nesciung eget distingy magndapibus.</p></td><td class="price" data-title="Price">'+old_price+'<span>$'+ob.price+'</span></td><td class="qty" data-title="Qty"><div class="input-group"><div class="button minus"> <button type="button" class="btn btn-primary btn-number" data-type="minus" data-field="quant['+ob.id+']"> <i class="ti-minus"></i> </button></div> <input type="text" name="quant['+ob.id+']" class="input-number" data-min="1" data-max="'+ob.stocks+'" value="'+ob.qty+'" data-pid="'+ob.id+'"><div class="button plus"> <button type="button" class="btn btn-primary btn-number" data-type="plus" data-field="quant['+ob.id+']"> <i class="ti-plus"></i> </button></div></div></td><td class="total-amount" data-title="Total"><span>$'+round(ob.qty*ob.price)+'</span></td><td class="action" data-title="Remove"><a href="#" onclick="removeFromCart('+ob.id+')"><i class="ti-trash remove-icon"></i></a></td></tr>'
				}
				html += '</tbody></table>';
				mainCart.html(html);

				$('#total-amount').html('<div class="row"><div class="col-lg-8 col-md-5 col-12"><div class="left"><div class="coupon"><form action="#" target="_blank"> <input name="Coupon" placeholder="Enter Your Coupon"> <button class="btn">Apply</button></form></div><div class="checkbox"> <label class="checkbox-inline" for="2"><input name="news" id="2" type="checkbox"> Gift Wrap (+10$)</label></div></div></div><div class="col-lg-4 col-md-7 col-12"><div class="right"><ul><li>Cart Subtotal<span>$'+round(res.payment.total)+'</span></li><li>Shipping<span>Free</span></li><li>You Save<span>$0.00</span></li><li class="last">You Pay<span>$'+round(res.payment.total)+'</span></li></ul><div class="button5"> <a href="./checkout.html" class="btn">Checkout</a> <a href="./home.html" class="btn">Continue shopping</a></div></div></div></div>')

				bindBtnCart(true)
			}

			const subTotal = $('#sub-total');
			if(subTotal.length) {
				subTotal.text('$'+round(res.payment.total));
				$('#total').text('$'+round(res.payment.total))
			}
		}
	})
  }

  function removeFromCart(pid){
	  AJAX("cart/"+cartId+"/"+pid, "DELETE", "", function(res, status, xhr){
		  if(xhr.status === 200) {
			  loadCartItems()
			  alert
		  }
	  })
  }

  function addToCart(pid, qty, type) {
	  if (!type) {
		  type = "POST"
	  }
	  const data = '{"id": '+pid+', "qty": '+qty+'}';
	  AJAX("cart/"+cartId, type, data, function(res, status, xhr){
		  if (xhr.status === 200) {
			  loadCartItems()
		  } else {
			  alert("Couldn't add to cart")
		  }
	  })
  }
