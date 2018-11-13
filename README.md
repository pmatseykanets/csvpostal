# csvpostal

A tiny utility that parses address info into components using [libpostal](https://github.com/openvenues/libpostal)

## Usage

By default it reads from `stdin` and writes to `stdout`

```bash
echo 'address\n"1 Main St., Baltimore, MD 21029, US"' | ./csvpostal
```

```txt
address,house,category,near,house_number,road,unit,level,staircase,entrance,po_box,postcode,suburb,city_district,city,island,state_district,state,country_region,country,world_region
"1 Main St., Baltimore, MD 21029, US",,,,1,main st.,,,,,,21029,,,baltimore,,,md,,us,
```

File in, file out

```bash
./csvpostal -in addresses.csv -out parsed_addresses.csv
```
