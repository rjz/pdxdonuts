window.makeMap = function(opts) {

  mapboxgl.accessToken = opts.accessToken;

  var map = new mapboxgl.Map({
    container: 'map',
    style: 'mapbox://styles/mapbox/light-v9',
    center: [opts.lng, opts.lat],
    zoom: 12.5,
  });

  var markers = [];

  // Global lock
  var activeMarker;

  var descriptionEl = document.getElementById('description');

  function clearDescription() {
    descriptionEl.className = '';
    descriptionEl.innerHTML = '';
  }

  function setDescription(title, desc) {
    descriptionEl.innerHTML = [
      '<h2>' + title + '</h2>',
      '<p>' + desc + '</p>',
    ].join('\n');

    setTimeout(function () {
      descriptionEl.className = 'active';
    });
  }

  function createArrows() {
    var arrows = [];
    var seed = Math.random() * 120;
    for (i = 0; i < 3; i++) {
      var rot = Math.floor((seed + Math.random() + i * 120) % 360);
      var scale = Math.random() * 2 + 1.3;
      var transform = [
        'scale(' + scale + ')',
        'rotate(' + rot + 'deg)',
        'translateY(' + Math.floor(-33) + 'px)',
      ].join(' ');
      arrows.push('<div class="arrow" style="transform:' + transform + '"></div>');
    }
    return arrows;
  }

  function MyMarker(place) {
    var iconSize = 64;
    var opts = {};
    var el = document.createElement('div');

    this._place = place;

    el.className = 'marker';
    el.innerHTML = [
      '<div class="shadow"></div>',
      '<div class="icon"></div>',
    ].concat(createArrows()).join('');

    mapboxgl.Marker.call(this, el, opts);
  }

  MyMarker.prototype = Object.create(mapboxgl.Marker.prototype);

  MyMarker.prototype._onMapClick = function (e) {
    var targetElement = e.originalEvent.target;
    var el = this._element;
    if (targetElement === el || el.contains(targetElement)) {
      el.className = 'marker active';
      map.panTo(this.getLngLat());

      // this.togglePopup();

      clearDescription();

      if (this === activeMarker) {
        this._element.className = 'marker';
        activeMarker = null;
        return;
      } else if (activeMarker) {
        activeMarker._element.className = 'marker';
        // activeMarker._popup.remove();
      }

      var properties = this._place.properties;

      setDescription(properties.name, properties.vicinity);
      activeMarker = this;
    }
  };

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
    // var popupOptions = {
    //   closeButton: false,
    //   closeOnClick: false,
    //   anchor: 'bottom',
    //   offset: {
    //     bottom: [0, -24]
    //   }
    // };

    // var detail = document.createElement('div');
    // detail.innerHTML = [
    //   '<h2>' + place.properties.name + '</h2>',
    //   '<p>' + place.properties.vicinity + '</p>'
    // ].join('');

    // var popup = new mapboxgl.Popup(popupOptions);
    // popup
    //   .setLngLat(place.geometry.coordinates)
    //   .setDOMContent(detail);

    var marker = new MyMarker(place);

    marker.setLngLat(place.geometry.coordinates)
    //  .setPopup(popup)
      .addTo(map);

    // popup.on('open', function () {
    //   setTimeout(function () {
    //     popup._content.className = 'mapboxgl-popup-content open';
    //   });
    // });

    // popup.on('close', function () {
    //   popup._content.className = 'mapboxgl-popup-content';
    // });

    markers.push(marker);
  });
};
