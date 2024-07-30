var sb_menu = [
	{
		id : "idSb1", icon: "mdi mdi-view-dashboard" , value: "HPC User Manage", data: [
			{ id: "idSb1-1", value: "Create User"},
			{ id: "idSb1-2", value: "Update User"},
			{ id: "idSb1-3", value: "Manage User"}
		]
	},
	{
		id : "idSb2", icon: "mdi mdi-table" , value: "HPC Job Manage", data: [
			{ id: "idSb2-1", value: "Job Status"},
			{ id: "idSb2-2", value: "Job Cancel"},
		]
	}
]
var v1 = {
	view: "form", id: "view1",
	rows: [
		{
			view: "datatable",
			columns: [
				{ "id": "title", "header": "Title", "fillspace": true, "sort": "string" },
				{ "id": "year", "header": "Year", "sort": "string" },
				{ "id": "votes", "header": "Votes", "sort": "string" },
				{ "id": "rating", "header": "Rating", "sort": "string" },
				{ "id": "rank", "header": "Rank", "sort": "string" },
				{ "id": "category", "header": "Category", "sort": "string" }
			],
			select: true,
			scrollX: false,
			url: ""
		},
		{
			view: "form",
			minHeight: 380,
			autoheight: false,
			elements: [
				{ view: "multicombo", "label": "To", "value": "2,3", "options": "" },
				{ view: "text", "label": "Subject" },
				{ view: "textarea", "label": "Message", "height": 150 },
				{ view: "button", "value": "Send Message", "css": "webix_success", "align": "center", "inputWidth": 200 }
			]
		}
	]
}
var v2 = {
	view: "form", id: "view2",
	rows: [
		{
			view: "datatable",
			columns: [
				{ "id": "title", "header": "Title", "fillspace": true, "sort": "string" },
				{ "id": "year", "header": "Year", "sort": "string" },
				{ "id": "votes", "header": "Votes", "sort": "string" },
				{ "id": "rating", "header": "Rating", "sort": "string" },
				{ "id": "rank", "header": "Rank", "sort": "string" },
				{ "id": "category", "header": "Category", "sort": "string" }
			],
			select: true,
			scrollX: false,
			url: ""
		},
		{
			view: "form",
			minHeight: 380,
			autoheight: false,
			elements: [
				{ view: "multicombo", "label": "To", "value": "2,3", "options": "" },
				{ view: "text", "label": "Subject" },
				{ view: "textarea", "label": "Message", "height": 150 },
				{ view: "button", "value": "Send Message", "css": "webix_primary", "align": "center", "inputWidth": 200 }
			]
		}
	]
}

webix.ready(function () {
    var layout = webix.ui({
	rows: [
		{ view: "toolbar", id: "idToolbar","css": "webix_dark", "padding": { "right": 10, "left": 10 },
			"elements": [
				{ view: "label", "label": "HPC Admin" }
			]
		},
		{
			id: "idColumn",
			type: "wide",
			cols: [
				{ view: "sidebar", id: "idSidebar" ,data:  sb_menu, "width": 200, 
					on: {
						onAfterSelect: function (id) {
							webix.message( this.getItem(id).value)
							var pos = $$("idColumn").index($$("idSidebar"));

							if ( this.getItem(id).value == "Job Status") {
								webix.ui(v2, $$("view1"));
								//$$("view1").disable();	
							}
							else
							{
								if ( $$("view2") != undefined) {
									webix.ui(v1, $$("view2"));
									//$$("view2").disable();	
								}
							}
							//("view1").hide();		
							//$$("view1").destructor();	
	
							//$$("view2").show();			
							let nview = $$("idColumn").getChildViews();
							 webix.message(nview)
							//$$("idColumn").addView(v2, 1)
						}
					}	
				},
				v1
			]
		}
	]
})
});

//$$("view2").hide();