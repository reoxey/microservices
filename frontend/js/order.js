
AJAX("orders", "GET", "", function(res, status, xhr){
    if(xhr.status === 200){
        const ordersList = $('#orders-list');

        const orderStatus = ["Ordered", "Order Paid", "Order Canceled", "Order Fulfilled", "Order Returned"];
        const PaymentMethod = ["Card", "COD", "PayPal"];

        let html = '';
        for(let ob of res){
            let textColour = "text-warning";
            switch (ob.status) {
                case 1: textColour = "text-success"; break;
                case 2: textColour = "text-danger"; break;
                case 3: textColour = "text-primary"; break;
                case 4: textColour = ""
            }
            html += '<section class="shop-services section home mt-1"><div class="container"><div class="row"><div class="col-6"><div class="single-service"><h3 class="mb-2">#'+ob.id+'</h4><h6>'+ob.created_at+'</p></div></div><div class="col-6"><div class="single-service"><h4 class="mb-2 text-right d-block w-100 '+textColour+'">'+orderStatus[ob.status]+'</h4><h6 class="pull-right text-right w-100 d-block">'+PaymentMethod[ob.payment.type]+'</h6></div></div><div class="col-2 mt-3"></div><div class="col-10 mt-3"><table class="table shopping-summery"><tbody>';

            for(let it of ob.items) {
                html += '<tr><td class="image" data-title="No"><img src="https://via.placeholder.com/50x50" alt="#"></td><td class="product-des" data-title="Description"><p class="product-name"><a href="#">'+it.name+'</a></p><p class="product-des">Maboriosam in a tonto nesciung eget distingy magndapibus.</p></td><td class="price" data-title="Price"><span>$'+it.price+'</span></td><td class="qty" data-title="Qty"> <span>'+it.qty+'</span></td><td class="total-amount" data-title="Total"><span>$'+round(it.qty*it.price)+'</span></td></tr>'
            }

            html += '<tr><td></td><td></td><td></td><td><strong>Total</strong></td><td class="total-amount" data-title="Total"><strong>$'+round(ob.payment.total)+'</strong></td></tr></tbody></table></div></div></div> </section>';
        }

        ordersList.html(html)
    }
})