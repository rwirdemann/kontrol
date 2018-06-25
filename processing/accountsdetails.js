Vue.component('my-accountdetails-bonus', {
  data: function () {
    console.log ("in vue component my-accountdetails-bonus: ", accountdetails.link);
    let sum = sumWithFilters(accountdetails.rows, []);
    let rueckstellung   = sumWithFilters(accountdetails.rows,
      [
        {"title" : "Type", "value": "Rückstellung"}
      ]
    );
    let costs   = sumWithFilters(accountdetails.rows,
      [
        {"title" : "Type", "value": "Gehalt"},
        {"title" : "Type", "value": "LNSteuer"},
        {"title" : "Type", "value": "SV-Beitrag"}
      ]
    );
    let revenue = sum-costs - rueckstellung ;
    let sales_p = sumWithFilters(accountdetails.rows, [{"title" : "Type", "value": "Vertriebsprovision"}]);
    let revenue07 = Math.round(70*(revenue-sales_p))/100;
    let bonus = revenue07 + costs + sales_p;
    let kshare = Math.round(100*(revenue + costs - sales_p - bonus))/100;
    return {
      accountdetails:  [
        {"key": "revenue", "value" : (revenue - sales_p)},
        {"key": "sales provision", "value" : sales_p},
        {"key": "costs", "value" : costs},
        {"key": "70% of revenue", "value" : revenue07},
        {"key": "Bonus", "value" : bonus},
        {"key": "Rückstellung", "value" : rueckstellung},
        {"key": "K-share", "value" : kshare}
      ]
    }
  },
  template: '<div><br>\
  <h3>bonus: Rule: employee gets 70% of his/her revenue minus costs plus sales bonus</h3>\
  <div v-for="pos in accountdetails">\
  <span>{{pos.key}}: </span><span>{{pos.value}}€</span>\
  </div>\
  <br></div>'
})

Vue.component('my-setfilter', {
  props: ['value'],
  data : function () {
    return {
      "foodata" : "",
      "deletefilter" : function (d) {return accountdetails.deleteFilter(d) },
      "setfilter" : function () {return accountdetails.setFilterText(this.foodata) },
      "filters" : accountdetails.filters };
    },

    "template": `<div>\
    <input v-model:value="foodata" placeholder="type filter">\
    <input type="submit" value="set filter" v-on:click="setfilter()">\
    <span v-if="filters[0]"><br>filters: </span>\
    <span class="filter"\
    v-for="filter in filters"\
    v-on:click="deletefilter(filter)"><span>{{filter.title}}: {{filter.value}}</span></span>\
    </div>`
  })

  function sumWithFilters (array, filters, s) {
    let sumOver = s || "Amount";
    let retval = 0;
    array.forEach (
      function(row) {
        let sumitup = false;
        if (filters.length == 0) {
          sumitup = true;
        }
        for (i in filters) {
          let d = filters[i];
          if (row[d.title] == d.value) sumitup = true;
        }
        if (sumitup) {
          retval += row[sumOver];
        }
      }
    );
    return Math.round(100*retval)/100;
  }


  function fetchAccountsDetails (link, domTarget, csFilter) {
    console.log ("in fetchAccountsDetails, link: ", link, domTarget);
    //

    //
    const fetchAsync = async () =>
    await (await fetch(link)).json()
    fetchAsync()
    .then( function (jsonData) {

      // define all the componets needed...
      console.log ("in fetchAsync: ", link, domTarget);

      Vue.component('my-tableheader', {
        props: ['cell'],
        template: '<div>{{cell}}</div>'
      }
    )
    Vue.component('my-tablebody', {
      props: ['row'],
      template: '<div>{{row}}</div>'
    }
  )

  // create the new vue
  var accountdetails = new Vue({
    el: '#'+domTarget,
    data: {
      rows: [],  //
      sortKey: '#',
      sortOrder: 'asc',
      link: link,
      domTarget: domTarget,
      test: "in fetchAccountsDetails, test",
      render: true,
      renderbonus: false,
      title: "this should be filled with different content,..",
      filters: []
    },
    methods: {
      "sum" : function sum (title) {
        mySum = 0;
        this.rows.forEach (function (row) {
          mySum += row[title] ;
        });
        return Math.round(100*mySum)/100;
      },
      "sortBy": function(sortKey) {
        this.reverse = (this.sortKey == sortKey) ? ! this.reverse : false;
        this.sortKey = sortKey;
      },
      "setFilterText": function setFilterText(str) {
        let filter = {};
        filter.title = "pattern";
        filter.value = str;
        if ( !this.filters.find(x => x.value  === str) ){
          this.filters.push (filter);
          this.rows = this.executeFilter();
        }
      },
      "setFilter": function setFilter(row, col) {
        let filter = {};
        filter.title = col;
        filter.value = row[col];
        if ( !this.filters.find(x => x.title === col) ) {
          this.filters.push (filter);
          this.rows = this.executeFilter();
        }
      },
      "deleteFilter": function deleteFilter(filter) {
        let index = this.filters.indexOf(filter);
        if ( index > -1 ) {
          this.filters.splice(index, 1);
        }
        this.rows = this.executeFilter();
      },
      "executeFilter": function executeFilter () {
        return this.flattenData.filter(
          function(row) {
            retval = true;
            this.forEach( function (d, i) {
              if (d.title == "pattern") {
                var cmpstr = "";
                for (var cell in row) {
                  cmpstr += row[cell].toString().toLowerCase();
                }
                if ( !cmpstr.includes(d.value.toLowerCase()) ) {
                  retval = false;
                }

              } else if (row[d.title] != d.value) {
                retval = false;
              }
            }
          );
          return retval;
        }, this.filters);
      },
      "logArray": function logArray (arr){
        console.log ("in logArray", arr);
      },
      "sortArray": function sortArray(array, sortKey) {
        let key = sortKey || this.sortKey || '#';
        let sortOption;
        if (this.sortOrder == 'asc') {
          this.sortOrder = 'dsc';
          sortOption = 1;
        } else {
          this.sortOrder = 'asc';
          sortOption = -1;
        }
        return array.sort(function(a, b) {
          if (a[key] > b[key]) {
            return 1*sortOption;
          } else if (a[key] < b[key]) {
            return  -1*sortOption;
          }
          return 0;
        }
      )
    },
  },
  computed: {
    "flattenData": function flattenData () {
      let retarray = [];
      console.log ("in flattenData: ", link.split(/[/ ]+/).pop() );
      if (link.split(/[/ ]+/).pop() == 'bankaccount') {
        // if link is like http://kommitment.dyn.amicdns.de:8991/kontrol/accounts
        for (let i in jsonData.Bookings ) {
          console.log ("in flattenData: --> bankaccount");
          retarray[i] = {};
          retarray[i]["Id"] = i;
          retarray[i]["Type"] = jsonData.Bookings[i].Type;
          retarray[i]["CostCenter"] = jsonData.Bookings[i].CostCenter;
          retarray[i]["Amount"] = Math.round(100*jsonData.Bookings[i].Amount)/100;
          retarray[i]["Year"] = jsonData.Bookings[i].Year;
          retarray[i]["Month"] = jsonData.Bookings[i].Month;
          retarray[i]["FileCreated"] = jsonData.Bookings[i].FileCreated;
          retarray[i]["BankCreated"] = jsonData.Bookings[i].BankCreated;
          retarray[i]["Text"] = jsonData.Bookings[i].Text;
        }
        this.title = "bankaccount"
      }
      else if (link.split(/[/ ]+/).pop() == 'accounts') {
        console.log ("in flattenData: --> accounts");
        this.title = "accounts";
        for (let i in jsonData.Accounts) {
          retarray[i] = {};
          retarray[i]["#"] = i;
          retarray[i]["Id"] = jsonData.Accounts[i].Owner.Id;
          retarray[i]["Name"] = jsonData.Accounts[i].Owner.Name;
          retarray[i]["Type"] = jsonData.Accounts[i].Owner.Type;
          retarray[i]["Revenue"] = Math.round(100*jsonData.Accounts[i].Revenue)/100;
          retarray[i]["Advances"] = Math.round(100*jsonData.Accounts[i].Advances)/100;
          retarray[i]["Internals"] = Math.round(100*jsonData.Accounts[i].Internals)/100;
          retarray[i]["sales commission"] = Math.round(100*jsonData.Accounts[i].Provision)/100;
          retarray[i]["Taxes"] = Math.round(100*jsonData.Accounts[i].Taxes)/100;
          retarray[i]["Costs"] = Math.round(100*jsonData.Accounts[i].Costs)/100;
          retarray[i]["Saldo"] = Math.round(100*jsonData.Accounts[i].Saldo)/100;
        }
      }
      else {
        // or http://kommitment.dyn.amicdns.de:8991/kontrol/accounts/AN
        for (let i in jsonData.Bookings ) {
          retarray[i] = {};
          retarray[i]["Id"] = i;
          retarray[i]["Type"] = jsonData.Bookings[i].Type;
          retarray[i]["CostCenter"] = jsonData.Bookings[i].CostCenter;
          retarray[i]["Amount"] = Math.round(100*jsonData.Bookings[i].Amount)/100;
          retarray[i]["Year"] = jsonData.Bookings[i].Year;
          retarray[i]["Month"] = jsonData.Bookings[i].Month;
          retarray[i]["FileCreated"] = jsonData.Bookings[i].FileCreated;
          retarray[i]["BankCreated"] = jsonData.Bookings[i].BankCreated;
          retarray[i]["Text"] = jsonData.Bookings[i].Text;
        }
        // fill the title
        this.title = jsonData.Owner.Name + ": " + jsonData.Owner.Type;
        // if ?cs filter is there then add it to the title
        if (csFilter) {
          this.title += " CostCenter="+csFilter;
        }
        console.log ("in flatten*, this.title = ",this.title);
      }
      return retarray;
    },
    "columns": function columns() {
      if (this.rows.length == 0) {
        return [];
      }
      return Object.keys(this.rows[0])
    }
  },
  created: function () {
    this.rows = this.flattenData;
    window.accountdetails = this;
    console.log ("... created fetchAccountsDetails", this.rows[0]);
    console.log ("... created fetchAccountsDetails", jsonData);
    // global.errors.push("cookies: \n", getCookies());
  }
});
})
.catch(reason => {
  console.log("fetchAccountsDetails: ",reason.message);
  global.errors.push("hi, something went wrong while conneting to the backend  ==> contact Johannes");
  global.errors.push("backend link: <a href='" +link+"'>"+link+"</a>");
})
}
