npm run build
rsync -avz --delete dist/ kr:/var/www/html/bluespot/
# rsync -avz --delete dist/ qyhever:/usr/share/nginx/html/bluespot/