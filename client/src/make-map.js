const { markerFactory } = require('./markers');

module.exports = function(opts) {
  mapboxgl.accessToken = opts.accessToken;

  const map = new mapboxgl.Map({
    container: 'map',
    style: 'mapbox://styles/mapbox/light-v9',
    center: [opts.lng, opts.lat],
    zoom: 12.5,
  });

  const createMarker = markerFactory(map);

  const dataSource = {
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

  const markers = dataSource.features.map(place => createMarker(place));

  // Jump through some hoops to get singleton tooltips working with custom icons
  // in MapboxGL.
  var activeMarkers = [];

  map.on('click', () => {
    window.setTimeout(() => {
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
};
