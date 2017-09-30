{{ define "scripts" }}
// This is a quick and dirty submitter.
(function($){
  var $downloadForm = $('.download-form');
  var $downloadUrls = $('.download-urls');
  var $resultField = $('.result-field');
  var csrfToken = $('input[name="gorilla.csrf.Token"]').val();

  // postJSON allows to post JSON serialized data to a given URL
  // while setting the CSRF header
  function postJSON(url, data) {
    return $.ajax({
      contentType: 'application/json',
      url: url,
      type: 'POST',
      headers: { 'X-CSRF-Token': csrfToken },
      data: JSON.stringify(data),
      dataType: 'json'
    });
  }

  $('.download-form').on('submit', function(e) {
    e.preventDefault();

    postJSON('/submit', { urls: $downloadUrls.val() }).done(function(response) {
      var linksCount = response.count;
      $resultField.attr({ 'class': 'alert alert-success' })
        .text(linksCount + ' ' + (linksCount > 1 ? 'links' : 'link') + ' processed.');
      $downloadUrls.val('');
    }).fail(function(xhr, textStatus) {
      $resultField.attr({ 'class': 'alert alert-danger' })
        .text('Error while processing given links: ' + xhr.responseText);
    });
  });
})(jQuery);
{{ end }}
