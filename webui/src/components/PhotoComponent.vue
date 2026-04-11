<template>
  <div class="photo-card glass-morphism">
    
    <div class="card-top">
      <div v-if="!IsOwner" class="user-info" @click="photoOwnerPressed">
        <div class="avatar-purple">{{ PhotoOwner.charAt(0).toUpperCase() }}</div>
        <span class="owner-name">@{{ PhotoOwner }}</span>
      </div>
      <div v-else class="delete-area">
        <button class="icon-btn delete" @click="deletePhoto" title="Elimina foto">
          <i class="fa-regular fa-trash-can"></i>
        </button>
      </div>
    </div>

    <div class="image-wrapper" @dblclick="heartButtonClicked">
      <img :src="PhotoUrl" class="main-photo" alt="Photo" data-bs-toggle="modal" :data-bs-target="'#PhotoFullSizeModal'+PhotoId">
    </div>

    <div class="interaction-bar">
      <div class="left-actions">
        <button class="btn-action" :class="{ 'is-liked': Liked }" @click="heartButtonClicked">
          <i :class="Liked ? 'fa-solid fa-heart' : 'fa-regular fa-heart'"></i>
          <span class="stat-num" data-bs-toggle="modal" :data-bs-target="'#LikesModal'+PhotoId" @click.stop="loadLikesList">
            {{ LikesNumber }}
          </span>
        </button>

        <button class="btn-action" data-bs-toggle="modal" :data-bs-target="'#CommentsModal'+PhotoId" @click="loadComments">
          <i class="fa-regular fa-comment"></i>
          <span class="stat-num">{{ CommentsNumber }}</span>
        </button>
      </div>
      
      <span class="upload-date">{{ formatDate(DateOfUpload) }}</span>
    </div>

    <div class="quick-comment">
      <input type="text" v-model="CommentText" placeholder="Aggiungi un commento..." @keyup.enter="commentPhoto">
      <button :disabled="CommentText.length == 0" @click="commentPhoto">Invia</button>
    </div>

    <div class="modal fade" :id="'LikesModal'+PhotoId" tabindex="-1">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content glass-morphism">
          <div class="modal-header border-0">
            <h5 class="modal-title text-white">Piace a</h5>
            <button type="button" class="btn-close btn-close-white" data-bs-dismiss="modal"></button>
          </div>
          <div class="modal-body">
            <div v-if="Likes.length === 0" class="text-muted">Ancora nessun like</div>
            <ul class="list-unstyled">
              <li v-for="user in Likes" :key="user" class="py-2 user-list-item" @click="goToUser(user)" data-bs-dismiss="modal">
                <div class="avatar-purple-sm">{{ user.charAt(0).toUpperCase() }}</div> @{{ user }}
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <div class="modal fade" :id="'CommentsModal'+PhotoId" tabindex="-1">
      <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable">
        <div class="modal-content glass-morphism">
          <div class="modal-header border-0">
            <h5 class="modal-title text-white">Commenti ({{ CommentsNumber }})</h5>
            <button type="button" class="btn-close btn-close-white" data-bs-dismiss="modal"></button>
          </div>
          <div class="modal-body">
            <div v-if="Comments.length === 0" class="text-muted">Ancora nessun commento</div>
            <div v-for="c in Comments" :key="c.CommentId" class="comment-item">
              <strong @click="goToUser(c.CommentAuthor)">@{{ c.CommentAuthor }}:</strong>
              <p>{{ c.CommentText }}</p>
              <hr class="op-1">
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="modal fade" :id="'PhotoFullSizeModal'+PhotoId" tabindex="-1">
      <div class="modal-dialog modal-xl modal-dialog-centered">
        <div class="modal-content bg-transparent border-0">
          <div class="modal-body p-0 text-center">
            <img :src="PhotoUrl" class="img-fluid rounded shadow-lg">
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
export default {
  name: 'PhotoComponent',
  props: ["owner", "photoid", "likesnumber", "commentsnumber", "dateofupload", "islikedbyme"],
  emits: ["photo_deleted_from_database"],
  
  data() {
    return {
      PhotoUrl: "",
      PhotoId: "",
      PhotoOwner: "",
      LikesNumber: 0,
      CommentsNumber: 0,
      Liked: false,
      Likes: [],
      Comments: [],
      CommentText: "",
      IsOwner: false
    }
  },

  methods: {
    formatDate(date) {
      if(!date) return "";
      return new Date(date).toLocaleDateString('it-IT', { day: 'numeric', month: 'short' });
    },

    async heartButtonClicked() {
      const loggedUser = localStorage.getItem("Username");
      const token = localStorage.getItem("SessionToken");
      const previousState = this.Liked;

      this.Liked = !this.Liked;
      this.LikesNumber += this.Liked ? 1 : -1;

      try {
        const config = { headers: { "Authorization": `Bearer ${token}` } };
        if (!previousState) {
          await this.$axios.post(`/users/${this.PhotoOwner}/photos/${this.PhotoId}/likes/`, {}, config);
        } else {
          await this.$axios.delete(`/users/${this.PhotoOwner}/photos/${this.PhotoId}/likes/${loggedUser}`, config);
        }
      } catch (e) {
        this.Liked = previousState;
        this.LikesNumber += this.Liked ? 1 : -1;
        console.error("Like Action Failed", e);
      }
    },

    async loadComments() {
      try {
        const res = await this.$axios.get(`/users/${this.PhotoOwner}/photos/${this.PhotoId}/comments/`);
        this.Comments = res.data.Comments || [];
      } catch (e) { console.error("Comments fetch error", e); }
    },

    async loadLikesList() {
      try {
        const res = await this.$axios.get(`/users/${this.PhotoOwner}/photos/${this.PhotoId}/likes/`);
        this.Likes = res.data.Likes ? res.data.Likes.map(l => l.Username) : [];
      } catch (e) { console.error("Likes fetch error", e); }
    },

    async commentPhoto() {
      const config = { headers: { "Authorization": `Bearer ${localStorage.getItem("Authstring")}` } };
      const payload = { 
        CommentAuthor: localStorage.getItem("Username"), 
        CommentText: this.CommentText 
      };
      try {
        const res = await this.$axios.post(`/users/${this.PhotoOwner}/photos/${this.PhotoId}/comments/`, payload, config);
        this.Comments.unshift(res.data);
        this.CommentsNumber++;
        this.CommentText = "";
      } catch (e) { console.error("Comment post error", e); }
    },

    async deletePhoto() {
      if (!confirm("Eliminare definitivamente questa foto?")) return;
      const config = { headers: { "Authorization": `Bearer ${localStorage.getItem("Authstring")}` } };
      try {
        await this.$axios.delete(`/users/${this.PhotoOwner}/photos/${this.PhotoId}/`, config);
        this.$emit("photo_deleted_from_database", this.PhotoId);
      } catch (e) { console.error("Delete failed", e); }
    },

    photoOwnerPressed() { this.$router.push(`/users/${this.PhotoOwner}/`); },
    goToUser(user) { this.$router.push(`/users/${user}/`); }
  },

  mounted() {
    this.PhotoId = this.photoid;
    this.PhotoOwner = this.owner;
    this.LikesNumber = parseInt(this.likesnumber || 0);
    this.CommentsNumber = parseInt(this.commentsnumber || 0);
    this.Liked = this.islikedbyme;
    this.IsOwner = (localStorage.getItem("Username") === this.PhotoOwner);
    
    // URL dinamico basato sulla configurazione axios
    const base = this.$axios.defaults.baseURL || "http://localhost:3000";
    this.PhotoUrl = `${base}/users/${this.owner}/photos/${this.photoid}/`;
  }
}
</script>

<style scoped>
.glass-morphism {
  background: rgba(45, 20, 60, 0.4);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(180, 100, 255, 0.2);
  border-radius: 18px;
  overflow: hidden;
  margin-bottom: 2rem;
}

.photo-card { width: 100%; max-width: 500px; color: white; transition: 0.3s; }
.card-top { padding: 12px 18px; display: flex; justify-content: space-between; align-items: center; }
.user-info { display: flex; align-items: center; cursor: pointer; gap: 10px; }
.avatar-purple {
  width: 32px; height: 32px; border-radius: 50%;
  background: linear-gradient(135deg, #8e2de2, #4a00e0);
  display: flex; align-items: center; justify-content: center; font-weight: bold; font-size: 0.8rem;
}

.image-wrapper { width: 100%; background: #000; display: flex; justify-content: center; overflow: hidden; }
.main-photo { max-width: 100%; height: auto; cursor: zoom-in; transition: transform 0.4s; }
.main-photo:hover { transform: scale(1.02); }

.interaction-bar { padding: 15px 18px; display: flex; justify-content: space-between; align-items: center; }
.btn-action { background: none; border: none; color: white; font-size: 1.3rem; cursor: pointer; transition: 0.2s; padding: 0 8px; }
.btn-action:hover { transform: scale(1.1); }
.is-liked { color: #ff2d55; text-shadow: 0 0 10px rgba(255, 45, 85, 0.5); }
.stat-num { font-size: 0.9rem; margin-left: 6px; font-weight: 300; }
.upload-date { font-size: 0.75rem; color: rgba(255,255,255,0.5); }

.quick-comment {
  border-top: 1px solid rgba(180, 100, 255, 0.1);
  padding: 12px 18px; display: flex; gap: 10px;
}
.quick-comment input {
  background: rgba(255,255,255,0.05); border: none; border-radius: 20px;
  padding: 5px 15px; color: white; flex-grow: 1; font-size: 0.9rem;
}
.quick-comment button { background: none; border: none; color: #b464ff; font-weight: bold; cursor: pointer; }
.icon-btn.delete:hover { color: #ff4757; }
.user-list-item { display: flex; align-items: center; gap: 10px; cursor: pointer; border-radius: 8px; padding: 5px; }
.user-list-item:hover { background: rgba(255,255,255,0.1); }
.avatar-purple-sm { width: 24px; height: 24px; border-radius: 50%; background: #8e2de2; font-size: 0.7rem; display: flex; align-items: center; justify-content: center; }
.op-1 { opacity: 0.1; color: white; }
</style>