import $ from 'jquery'
import Notify from 'handy-notification'

// FUNCTION FOR COMMON LOGIN
const commonLogin = options => {
  let { data, btn, url, redirect, defBtnValue } = options
  let overlay2 = $('.overlay-2')

  btn
    .attr('value', 'Please wait..')
    .addClass('a_disabled')
  overlay2.show()

  $.ajax({
    url,
    data,
    method: 'POST',
    dataType: 'JSON',
    success: (data) => {
      let { mssg, success } = data
      if (success) {
        Notify({ value: mssg, done: () => location.href = redirect })
        btn.attr('value', 'Redirecting..')
        overlay2.show()
      } else {
        Notify({ value: mssg })
        btn
          .attr('value', defBtnValue)
          .removeClass('a_disabled')
        overlay2.hide()
      }
      btn.blur()
    }
  })
}

// FUNCTION FOR GOING BACK
const Back = e => {
  e.preventDefault()
  window.history.back()
}

module.exports = {
  commonLogin,
  Back
}
