sudo -i -u postgres psql -d greenlight -c "GRANT ALL PRIVILEGES ON DATABASE greenlight TO greenlight"
sudo -i -u postgres psql -d greenlight -c "GRANT ALL ON SCHEMA public TO greenlight"