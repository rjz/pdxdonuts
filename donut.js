window.makeMap = function(opts) {

  mapboxgl.accessToken = opts.accessToken;

  var map = new mapboxgl.Map({
    container: 'map',
    style: 'mapbox://styles/mapbox/light-v9',
    center: [opts.lng, opts.lat],
    zoom: 12.5,
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

  dataSource.features.forEach(function (place) {
    var iconSize = 64;
    var el = document.createElement('div');
    el.className = 'marker';
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
  });
};
