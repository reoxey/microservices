
// AJAX("carts", "POST", "", function(res, status, xhr){
    
// })

AJAX("products", "GET", "", function(res, status, xhr){

    if(xhr.status === 200) {
        const obj = $("#all-products");
        for(let ob of res) {

            const dO = new Date(ob.created_at);
            const dN = new Date();

            let tag = '';
            if(dN.getTime()-dO.getTime() < 518500590) 
                tag = '<span class="new">New</span>';
            if (ob.stocks <= 0) tag = '<span class="out-of-stock">Out of stock</span>';
            else if (ob.stocks < 5) tag = '<span class="price-dec">'+ob.stocks+' left</span>';

            obj.append('<div class="col-lg-4 col-md-6 col-12">'+
            '<div class="single-product">'+
                '<div class="product-img">'+
                    '<a title="Quick View" href="#" onclick="modalData('+ob.id+','+ob.stocks+')">'+
                        '<img class="default-img" src="https://via.placeholder.com/550x750" alt="#">'+
                        '<img class="hover-img" src="https://via.placeholder.com/550x750" alt="#">'+
                        tag+
                    '</a>'+
                    '<div class="button-head">'+
                        '<div class="product-action">'+
                            '<a title="Quick View" href="#" onclick="modalData('+ob.id+','+ob.stocks+')"><i class=" ti-eye"></i><span>Quick Shop</span></a>'+
                            '<a title="Wishlist" href="#"><i class=" ti-heart "></i><span>Add to Wishlist</span></a>'+
                            '<a title="Compare" href="#"><i class="ti-bar-chart-alt"></i><span>Add to Compare</span></a>'+
                        '</div>'+
                        '<div class="product-action-2">'+
                            '<a title="Add to cart" href="#" onclick="addToCart('+ob.id+',1)">Add to cart</a>'+
                        '</div>'+
                    '</div>'+
                '</div>'+
                '<div class="product-content">'+
                    '<h3><a id="product-title-'+ob.id+'" title="Quick View" href="#" onclick="modalData('+ob.id+','+ob.stocks+')">'+ob.name+'</a></h3>'+
                    '<div class="product-price">'+
                        '<span id="product-price-'+ob.id+'">$'+ob.price+'</span>'+
                    '</div>'+
                '</div>'+
            '</div>'+
        '</div>')
        }
    }
})


