SET
SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET
time_zone = "+00:00";

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;


CREATE TABLE `carts`
(
    `id`         int(11) NOT NULL,
    `buyer`      int(11) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `carts_items`
(
    `id`        int(11) NOT NULL,
    `cart_id`   int(11) NOT NULL,
    `name`      varchar(100) NOT NULL,
    `price`     float        NOT NULL,
    `old_price` float        NOT NULL DEFAULT 0,
    `stocks`    int(5) NOT NULL,
    `qty`       int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `orders`
(
    `id`         int(11) NOT NULL,
    `cart_id`    int(11) NOT NULL,
    `buyer`      int(11) NOT NULL,
    `amount`     float     NOT NULL,
    `pay_method` tinyint(1) NOT NULL,
    `status`     tinyint(1) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `orders_items`
(
    `id`       int(11) NOT NULL,
    `order_id` int(11) NOT NULL,
    `name`     varchar(100) NOT NULL,
    `price`    float        NOT NULL,
    `qty`      int(4) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `products`
(
    `id`         int(11) NOT NULL,
    `sku`        varchar(7)  NOT NULL,
    `name`       varchar(50) NOT NULL,
    `price`      double      NOT NULL,
    `stocks`     int(5) NOT NULL DEFAULT 1,
    `created_at` timestamp   NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `shipping`
(
    `id`             int(11) NOT NULL,
    `order_id`       int(11) NOT NULL,
    `status`         tinyint(1) NOT NULL,
    `payment_handle` tinyint(1) NOT NULL,
    `payment_status` tinyint(1) NOT NULL,
    `user`           int(11) NOT NULL,
    `contact_name`   varchar(100) NOT NULL,
    `contact_phone`  varchar(100) NOT NULL,
    `landmark`       varchar(500) NOT NULL,
    `city`           varchar(100) NOT NULL,
    `state`          varchar(100) NOT NULL,
    `country`        varchar(100) NOT NULL,
    `zip`            int(6) NOT NULL,
    `created_at`     timestamp    NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `shipping_addresses`
(
    `id`            int(11) NOT NULL,
    `user`          int(11) NOT NULL,
    `contact_name`  varchar(100) NOT NULL,
    `contact_phone` varchar(100) NOT NULL,
    `landmark`      varchar(500) NOT NULL,
    `city`          varchar(100) NOT NULL,
    `state`         varchar(100) NOT NULL,
    `country`       varchar(100) NOT NULL,
    `zip`           int(6) NOT NULL,
    `created_at`    timestamp    NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `users`
(
    `id`        int(11) NOT NULL,
    `name`      varchar(50)  NOT NULL,
    `email`     varchar(100) NOT NULL,
    `password`  varchar(200) NOT NULL,
    `is_admin`  tinyint(1) NOT NULL DEFAULT 0,
    `joined_at` timestamp    NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


ALTER TABLE `carts`
    ADD PRIMARY KEY (`id`);

ALTER TABLE `carts_items`
    ADD PRIMARY KEY (`id`, `cart_id`);

ALTER TABLE `orders`
    ADD PRIMARY KEY (`id`) USING BTREE;

ALTER TABLE `orders_items`
    ADD PRIMARY KEY (`id`, `order_id`);

ALTER TABLE `products`
    ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `sku` (`sku`);

ALTER TABLE `shipping`
    ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `order_id` (`order_id`);

ALTER TABLE `shipping_addresses`
    ADD PRIMARY KEY (`id`);

ALTER TABLE `users`
    ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`);


ALTER TABLE `carts`
    MODIFY `id` int (11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `orders`
    MODIFY `id` int (11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `products`
    MODIFY `id` int (11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `shipping`
    MODIFY `id` int (11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `shipping_addresses`
    MODIFY `id` int (11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `users`
    MODIFY `id` int (11) NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
