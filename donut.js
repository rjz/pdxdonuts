window.makeMap = function(opts) {

  mapboxgl.accessToken = opts.accessToken;

  var map = new mapboxgl.Map({
    container: 'map',
    style: 'mapbox://styles/mapbox/light-v9',
    center: [opts.lng, opts.lat],
    zoom: 12.5,
  });

  var markers = [];

  // Jump through some hoops to get singleton tooltips working with custom icons
  // in MapboxGL.
  var activeMarkers = [];

  map.on('click', function () {
    window.setTimeout(function () {
      // Find all currently-open popups
      var ms = markers.filter(function (m) {
        return m.getPopup().isOpen();
      });

      // Close all *previously*-open popups
      ms.forEach(function (m) {
        if (activeMarkers.indexOf(m) > -1) {
          m.togglePopup();
        }
      });

      // Pan to currently-open popup (there *should* only be one at a time)
      if (ms.length) {
        map.panTo(ms[0].getLngLat());
      }

      // Update previously-open list
      activeMarkers = ms;
    });
  });

  var dataSource = {
    type: 'FeatureCollection',
    features: opts.data.map(function (place) {
      return {
        type: 'Feature',
        properties: {
          name: place.name,
          vicinity: place.vicinity,
        },
        geometry: {
          type: 'Point',
          coordinates: [
            place.location.lng,
            place.location.lat
          ],
        }
      }
    }),
  };

  var THEME_COUNT = 4;

  dataSource.features.forEach(function (place) {
    var iconSize = 64;
    var el = document.createElement('div');
    var index = place.properties.name.split('')
      .reduce((sum, c) => sum + c.charCodeAt(0), 0) % THEME_COUNT;
    el.className = 'marker marker-' + index;
    el.style.width = iconSize + 'px';
    el.style.height = iconSize + 'px';
    el.style.borderRadius = Math.floor(iconSize / 2) + 'px';

    var popupOptions = {
      closeButton: false,
      closeOnClick: false,
      anchor: 'bottom',
      offset: {
        bottom: [0, -28]
      }
    };

    var detail = document.createElement('div');
    detail.className = 'detail-' + index;
    detail.innerHTML = [
      '<h2>' + place.properties.name + '</h2>',
      '<p>' + place.properties.vicinity + '</p>'
    ].join('');

    var popup = new mapboxgl.Popup(popupOptions)
      .setLngLat(place.geometry.coordinates)
      .setDOMContent(detail)
      .addTo(map);

    var marker = new mapboxgl.Marker(el, {
      offset: [
        -iconSize / 2,
        -iconSize / 2
      ]
    });

    marker.setLngLat(place.geometry.coordinates)
      .setPopup(popup)
      .addTo(map);

    markers.push(marker);
  });
};
