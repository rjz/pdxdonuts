const THEME_COUNT = 4;

const ICON_SIZE = 64;

const checksum = str =>
  str.split('').reduce((sum, c) => sum + c.charCodeAt(0), 0) % THEME_COUNT;

const createMarkerElement = () => {
  const el = document.createElement('div');
  el.style.width = ICON_SIZE + 'px';
  el.style.height = ICON_SIZE + 'px';
  el.style.borderRadius = Math.floor(ICON_SIZE / 2) + 'px';
  return el;
};

const createDetailElement = ({ name, vicinity }) => {
  const el = document.createElement('div');
  el.innerHTML = [
    '<h2>' + name + '</h2>',
    '<p>' + vicinity + '</p>'
  ].join('');
  return el;
};

const POPUP_OPTIONS = {
  closeButton: false,
  closeOnClick: false,
  anchor: 'bottom',
  offset: {
    bottom: [0, -28]
  }
};

const MARKER_OPTIONS = {
  offset: [
    -ICON_SIZE / 2,
    -ICON_SIZE / 2
  ]
};

// TODO: enable overrides of constantized params
module.exports.markerFactory = map => (place) => {
  const index = checksum(place.properties.name);

  const detail = createDetailElement(place.properties);
  detail.className = 'el-' + index;

  const popup = new mapboxgl.Popup(POPUP_OPTIONS)
    .setLngLat(place.geometry.coordinates)
    .setDOMContent(detail)
    .addTo(map);

  const markerEl = createMarkerElement();
  markerEl.className = 'marker marker-' + index;

  const marker = new mapboxgl.Marker(markerEl, MARKER_OPTIONS);
  marker.setLngLat(place.geometry.coordinates)
    .setPopup(popup)
    .addTo(map);

  return marker;
};
