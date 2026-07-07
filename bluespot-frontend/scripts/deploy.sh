npm run build
rsync -avz --delete dist/ kr:/var/www/html/sfr/
# rsync -avz --delete dist/ qyhever:/usr/share/nginx/html/sfr/