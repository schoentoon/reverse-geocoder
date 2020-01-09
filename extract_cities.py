#!/usr/bin/python

"""
This code is pretty much copied from https://github.com/thampiman/reverse-geocoder/blob/master/reverse_geocoder/__init__.py
All it does is download the cities1000 file from geonames.org and merge admin1 and admin2 codes into it.
It made more sense to me to simply include this in this repository over linking to it.
"""

import os
import zipfile
import csv

"""
Below the files to download are specified, by default it'll use the cities1000 file which is as by description of geonames themselves:
cities1000.zip           : all cities with a population > 1000 or seats of adm div down to PPLA3 (ca 130.000), see 'geoname' table for columns
You can however freely change this to their 500/5000/15000 versions for higher accuracy or simply a smaller tree and possibly more performance.
"""

GN_URL = 'http://download.geonames.org/export/dump/'
GN_CITIES1000 = 'cities1000'
GN_ADMIN1 = 'admin1CodesASCII.txt'
GN_ADMIN2 = 'admin2Codes.txt'

# Schema of the GeoNames Cities with Population > 1000
GN_COLUMNS = {
    'geoNameId': 0,
    'name': 1,
    'asciiName': 2,
    'alternateNames': 3,
    'latitude': 4,
    'longitude': 5,
    'featureClass': 6,
    'featureCode': 7,
    'countryCode': 8,
    'cc2': 9,
    'admin1Code': 10,
    'admin2Code': 11,
    'admin3Code': 12,
    'admin4Code': 13,
    'population': 14,
    'elevation': 15,
    'dem': 16,
    'timezone': 17,
    'modificationDate': 18
}

# Schema of the GeoNames Admin 1/2 Codes
ADMIN_COLUMNS = {
    'concatCodes': 0,
    'name': 1,
    'asciiName': 2,
    'geoNameId': 3
}

# Schema of the cities file created by this library
RG_COLUMNS = [
    'lat',
    'lon',
    'name',
    'admin1',
    'admin2',
    'cc'
]

# Name of cities file created by this library
RG_FILE = 'rg_cities.csv'

gn_cities1000_url = GN_URL + GN_CITIES1000 + '.zip'
gn_admin1_url = GN_URL + GN_ADMIN1
gn_admin2_url = GN_URL + GN_ADMIN2

cities1000_zipfilename = GN_CITIES1000 + '.zip'
cities1000_filename = GN_CITIES1000 + '.txt'

if not os.path.exists(cities1000_zipfilename):
    print('Downloading files from Geoname...')
    try: # Python 3
        import urllib.request
        urllib.request.urlretrieve(gn_cities1000_url, cities1000_zipfilename)
        urllib.request.urlretrieve(gn_admin1_url, GN_ADMIN1)
        urllib.request.urlretrieve(gn_admin2_url, GN_ADMIN2)
    except ImportError: # Python 2
        import urllib
        urllib.urlretrieve(gn_cities1000_url, cities1000_zipfilename)
        urllib.urlretrieve(gn_admin1_url, GN_ADMIN1)
        urllib.urlretrieve(gn_admin2_url, GN_ADMIN2)

_z = zipfile.ZipFile(open(cities1000_zipfilename, 'rb'))
open(cities1000_filename, 'wb').write(_z.read(cities1000_filename))

print('Loading admin1 codes...')
admin1_map = {}
t_rows = csv.reader(open(GN_ADMIN1, 'rt'), delimiter='\t')
for row in t_rows:
    admin1_map[row[ADMIN_COLUMNS['concatCodes']]] = row[ADMIN_COLUMNS['asciiName']]

print('Loading admin2 codes...')
admin2_map = {}
for row in csv.reader(open(GN_ADMIN2, 'rt'), delimiter='\t'):
    admin2_map[row[ADMIN_COLUMNS['concatCodes']]] = row[ADMIN_COLUMNS['asciiName']]


print('Creating formatted geocoded file...')
writer = csv.DictWriter(open(RG_FILE, 'wt'), fieldnames=RG_COLUMNS)
rows = []
for row in csv.reader(open(cities1000_filename, 'rt'), \
        delimiter='\t', quoting=csv.QUOTE_NONE):
    lat = row[GN_COLUMNS['latitude']]
    lon = row[GN_COLUMNS['longitude']]
    name = row[GN_COLUMNS['asciiName']]
    cc = row[GN_COLUMNS['countryCode']]

    admin1_c = row[GN_COLUMNS['admin1Code']]
    admin2_c = row[GN_COLUMNS['admin2Code']]

    cc_admin1 = cc+'.'+admin1_c
    cc_admin2 = cc+'.'+admin1_c+'.'+admin2_c

    admin1 = ''
    admin2 = ''

    if cc_admin1 in admin1_map:
        admin1 = admin1_map[cc_admin1]
    if cc_admin2 in admin2_map:
        admin2 = admin2_map[cc_admin2]

    write_row = {'lat':lat,
                    'lon':lon,
                    'name':name,
                    'admin1':admin1,
                    'admin2':admin2,
                    'cc':cc}
    rows.append(write_row)
writer.writeheader()
writer.writerows(rows)

print('Cleaning up..')
os.remove(cities1000_filename)
os.remove(cities1000_zipfilename)
os.remove(GN_ADMIN1)
os.remove(GN_ADMIN2)