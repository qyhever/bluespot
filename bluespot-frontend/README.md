# bluespot-frontend

## common commands
```
rsync -avz --progress --partial ./scripts/nginx.conf kr:/etc/nginx/sites-available/bluespot.6688988.xyz.conf
ln -s /etc/nginx/sites-available/bluespot.6688988.xyz.conf /etc/nginx/sites-enabled/
nginx -t
nginx -s reload
sudo certbot --nginx -d bluespot.6688988.xyz
```
