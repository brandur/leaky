

    ADDR=:5100 PEER_URLS=http://localhost:5200 ./leaky
    ADDR=:5200 PEER_URLS=http://localhost:5100 ./leaky
